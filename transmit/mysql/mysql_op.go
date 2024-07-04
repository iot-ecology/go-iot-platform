package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

// 全局map，用于存储id和数据库连接的映射
var dbMap = make(map[uint]*sql.DB)

// 用于同步的互斥锁，防止并发访问map时发生冲突
var mu sync.Mutex

func InitMySQLConnection(username, host, password, dbname string, port int, id uint) (*sql.DB, error) {
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
