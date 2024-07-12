package cassandra

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"iot-transmit/common"
	"strings"
	"sync"
	"time"
)

var cassandraClientMap = make(map[uint]*gocql.Session) // 存储id和Cassandra会话的映射
var mu sync.Mutex                                      // 用于同步的互斥锁

func GetCassandra(ips []string, username, password string, id uint) (*gocql.Session, error) {
	mu.Lock()
	defer mu.Unlock()

	// 尝试获取现有的Cassandra会话
	session, exists := cassandraClientMap[id]
	if exists {
		zap.S().Infof("Reusing existing Cassandra session for id: %d", id)
		return session, nil
	}

	// 创建Cassandra集群配置
	cluster := gocql.NewCluster(ips...)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username:              username,
		Password:              password, // 请确保AllowedAuthenticators中的authenticator是您Cassandra集群支持的
		AllowedAuthenticators: []string{"com.instaclustr.cassandra.auth.InstaclustrPasswordAuthenticator"},
	}

	// 创建会话
	session, err := cluster.CreateSession()
	if err != nil {
		zap.S().Fatalf("error %+v", err)
		return nil, err
	}

	// 将新的Cassandra会话存储到map中
	cassandraClientMap[id] = session
	zap.S().Infof("New Cassandra session for id %d has been initiated.\n", id)

	return session, nil
}

type CassandraOp struct{}

func (op *CassandraOp) HandleDataRowLists(db, table, Script string, dataRowList []common.DataRowList, client *gocql.Session) error {
	res := op.RunScript(dataRowList, Script)

	var sqlStatements []string
	for _, re := range res {
		sqlStatements = append(sqlStatements, op.Save(re, table, db))
	}
	for _, statement := range sqlStatements {
		zap.S().Infof("执行SQL语句：%s", statement)
		err := client.Query(statement).Exec()
		if err != nil {
			zap.S().Errorf("执行SQL语句失败！")
		}
	}

	return nil
}

func (op *CassandraOp) Save(dt []CassandraParam, table string, db string) string {
	var fields []string
	var valueStrs []string
	for _, param := range dt {
		fields = append(fields, "\""+param.FieldName+"\"")        // 使用反引号包围字段名
		valueStrs = append(valueStrs, buildValueStr(param.Value)) // 使用自定义函数处理值的格式化
	}

	fieldStr := strings.Join(fields, ", ")
	valueStr := strings.Join(valueStrs, ", ")

	query := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s);", db, table, fieldStr, valueStr)
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

func (op *CassandraOp) RunScript(dataRowList []common.DataRowList, script string) [][]CassandraParam {

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题！")
		return nil
	}
	var fn func(string2 []common.DataRowList) [][]CassandraParam
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		zap.S().Errorf("Js函数映射到 Go 函数失败！")
		return nil
	}
	a := fn(dataRowList)

	return a

}

type CassandraParam struct {
	FieldName string      `json:"field_name"`
	Value     interface{} `json:"value"`
}
