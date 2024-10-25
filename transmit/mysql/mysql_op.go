package mysql

import (
	"database/sql"
	"fmt"
	"iot-transmit/common"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// 全局map，用于存储id和数据库连接的映射
var dbMap = make(map[string]*sql.DB)

// 用于同步的互斥锁，防止并发访问map时发生冲突
var mu sync.Mutex

func InitMySQLConnection(username, host, password, dbname string, port int, id string) (*sql.DB, error) {

	mu.Lock()         // 进入临界区前加锁
	defer mu.Unlock() // 确保在函数返回时释放锁

	// 首先尝试从map中获取现有的数据库连接
	if db, exists := dbMap[id]; exists {
		zap.S().Infof("Reusing existing MySQL connection for id %d.", id)
		return db, nil
	}

	// 如果map中没有这个id的连接，则创建新的连接
	dsn := username + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		zap.S().Fatalf("database open failed, err: %v", err)
		return nil, err
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(time.Minute * 60) // 设置连接的最大生命周期

	// 测试数据库连接
	err = DB.Ping()
	if err != nil {
		zap.S().Fatalf("database ping failed, err: %v", err)
		return nil, err
	}

	// 将新的数据库连接存储到map中
	dbMap[id] = DB
	zap.S().Infof("New MySQL connection pool for id %d has been initiated.", id)
	return DB, nil
}

type MysqlOp struct{}

func (op *MysqlOp) HandleDataRowLists(table, Script string, dataRowList []common.DataRowList, client *sql.DB) error {
	res := op.RunScript(dataRowList, Script)

	var sqlStatements []string
	for _, re := range res {
		sqlStatements = append(sqlStatements, op.Save(re, table))
	}

	tx, err := client.Begin()
	if err != nil {
		zap.S().Errorf("Begin transaction failed, err: %v", err)
	}

	for _, statement := range sqlStatements {
		_, err = tx.Exec(statement)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				zap.S().Errorf("Rollback transaction failed, err: %v", err)
			}
			zap.S().Errorf("Exec sql failed, err: %v", err)
		}
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		zap.S().Errorf("Commit transaction failed, err: %v", err)
	}

	return nil
}

func (op *MysqlOp) Save(dt []MysqlParam, table string) string {
	var fields []string
	var valueStrs []string
	for _, param := range dt {
		fields = append(fields, "`"+param.FieldName+"`")          // 使用反引号包围字段名
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

func (op *MysqlOp) RunScript(dataRowList []common.DataRowList, script string) [][]MysqlParam {

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题！")
		return nil
	}
	var fn func(string2 []common.DataRowList) [][]MysqlParam
	get := vm.Get("main")
	if get != nil {

		err = vm.ExportTo(get, &fn)
		if err != nil {
			zap.S().Errorf("Js函数映射到 Go 函数失败！")
			return nil
		}
		a := fn(dataRowList)
		return a
	}
	return nil
}

type MysqlParam struct {
	FieldName string
	Value     interface{}
}
