package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"iot-transmit/common"
	"net"
	"strings"
	"sync"
	"time"
)

// 全局map，用于存储id和ClickHouse连接的映射
var clickHouseClientMap = make(map[uint]driver.Conn)

// 用于同步的互斥锁
var mu sync.Mutex

// GetClickHouse 创建或获取ClickHouse连接
func GetClickHouse(id uint, addr []string, database, username, password string) (driver.Conn, error) {
	mu.Lock()
	defer mu.Unlock()

	session, exists := clickHouseClientMap[id]

	if exists {
		err := session.Ping(context.Background())
		if err != nil {
			// 如果连接出现问题，则移除并重新创建
			delete(clickHouseClientMap, id)
		} else {
			fmt.Printf("Reusing existing ClickHouse connection for id: %d.\n", id)
			return session, nil
		}
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: addr,
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	if err != nil {
		zap.S().Fatal("Failed to create ClickHouse connection", zap.Error(err))
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		zap.S().Fatal("Failed to create ClickHouse connection", zap.Error(err))
		return nil, err
	}
	clickHouseClientMap[id] = conn
	zap.S().Infof("New ClickHouse connection for id %d has been initiated", id)
	return conn, nil
}

type ClickhouseOp struct{}

func (op *ClickhouseOp) HandleDataRowLists(

	table, Script string, dataRowList []common.DataRowList, client driver.Conn) error {
	res := op.RunScript(dataRowList, Script)

	var sqlStatements []string
	for _, re := range res {
		sqlStatements = append(sqlStatements, op.Save(re, table))
	}

	for _, statement := range sqlStatements {
		err := client.Exec(context.Background(), statement)
		if err != nil {
			zap.S().Errorf("执行SQL语句失败：%s", statement)
		}
	}

	return nil
}

func (op *ClickhouseOp) Save(dt []ClickhouseParam, table string) string {
	var fields []string
	var valueStrs []string
	for _, param := range dt {
		fields = append(fields, "\""+param.FieldName+"\"")        // 使用反引号包围字段名
		valueStrs = append(valueStrs, buildValueStr(param.Value)) // 使用自定义函数处理值的格式化
	}

	fieldStr := strings.Join(fields, ", ")
	valueStr := strings.Join(valueStrs, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, fieldStr, valueStr)
	return query
}

func buildValueStr(value interface{}) string {
	switch v := value.(type) {
	case string:
		return "'" + v + "'"
	case int, int64, float64, bool:
		return fmt.Sprintf("%v", value)
	default:
		// 对于其他类型，这里可以添加更多的处理逻辑
		return "NULL"
	}
}

func (op *ClickhouseOp) RunScript(dataRowList []common.DataRowList, script string) [][]ClickhouseParam {

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题！")
		return nil
	}
	var fn func(string2 []common.DataRowList) [][]ClickhouseParam
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		zap.S().Errorf("Js函数映射到 Go 函数失败！")
		return nil
	}
	a := fn(dataRowList)
	return a

}

type ClickhouseParam struct {
	FieldName string
	Value     interface{}
}
