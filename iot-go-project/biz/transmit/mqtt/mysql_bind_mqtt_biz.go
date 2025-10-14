/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package transmit

import (
	"context"
	"encoding/json"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"iot-transmit/cache"
	"iot-transmit/common"
	"iot-transmit/mysql"
	"strconv"
)

type MySQLTransmitBindBiz struct{}

func (biz *MySQLTransmitBindBiz) PageData(table string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.MySQLTransmitBind

	db := glob.GDb

	if table != "" {
		db = db.Where("table like ?", "%"+table+"%")
	}

	db.Model(&models.MySQLTransmitBind{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}




func (biz *MySQLTransmitBindBiz) Bind(req models.MySQLTransmitBind) {
	if req.Enable {
		jsonData := biz.toByte(req)
		// 缓存构造
		glob.GRedis.LPush(context.Background(), "transmit:mysql:" + req.Protocol+":"+strconv.Itoa(req.DeviceUid) +":" +req.
			IdentificationCode,
			jsonData)
	}
}

func (biz *MySQLTransmitBindBiz) toByte(req models.MySQLTransmitBind) []byte {
	var mysqlInfo models.MySQLTransmit

	glob.GDb.First(&mysqlInfo, req.MySQLTransmitId)

	v := cache.MySQLTransmitCache{
		ID:       "mysql-" + strconv.Itoa(int(req.ID)),
		Host:     mysqlInfo.Host,
		Port:     mysqlInfo.Port,
		Username: mysqlInfo.Username,
		Password: mysqlInfo.Password,
		Database: mysqlInfo.Database,
		Table:    req.Table,
		Script:   req.Script,
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// ChangeEnable 修改启用状态
func (biz *MySQLTransmitBindBiz) ChangeEnable(req models.MySQLTransmitBind) {
	glob.GDb.Model(models.MySQLTransmitBind{}).Where("id = ?", req.ID).Update("enable", req.Enable)
	biz.HandlerRedis(req)
}

func (biz *MySQLTransmitBindBiz) HandlerRedis(req models.MySQLTransmitBind) {
	if req.Enable {
		biz.Bind(req)
	} else {

		glob.GRedis.LRem(context.Background(), "transmit:mysql:"+ req.Protocol +":"+strconv.Itoa(req.DeviceUid) +":" +req.
			IdentificationCode,	1,
			biz.toByte(req))
	}
}

var mysqlOp = mysql.MysqlOp{}

// MockScript 模拟执行脚本
func (biz *MySQLTransmitBindBiz) MockScript(dataRowList []common.DataRowList, script string) [][]mysql.MysqlParam {
	return mysqlOp.RunScript(dataRowList, script)
}
