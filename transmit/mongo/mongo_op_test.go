package mongo

import (
	"testing"
)

func TestGet(t *testing.T) {
	// 示例：获取或创建MongoDB客户端
	Host := "localhost"
	Username := "admin"
	Password := "admin"
	Db := "iot"
	Port := 27017
	client, err := GetMongoDBClient(Host, Username, Password, Db, Port, 1)
	if err != nil {
		t.Fatal(err)
	}
	client2, err2 := GetMongoDBClient(Host, Username, Password, Db, Port, 1)
	if err2 != nil {
		t.Fatal(err)
	}
	if client == client2 {
		t.Logf("success")
	}
}
