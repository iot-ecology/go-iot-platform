package main

import (
	"encoding/json"
	"go.uber.org/zap"
)

type CassandraHandler struct {
}

func (handler CassandraHandler) handler(body []byte) {
	var data []DataRowList
	err := json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("处理cassandra数据失败", zap.Error(err))
	}
}
