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

package biz

import (
	"fmt"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"igp/ut"
	"time"

	"github.com/dop251/goja"
)

type SignalDelayWaringBiz struct{}

func (biz *SignalDelayWaringBiz) PageData(name string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dt []models.SignalDelayWaring

	db := glob.GDb

	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	db.Model(&models.SignalDelayWaring{}).Count(&pagination.Total)

	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dt)

	pagination.Data = dt
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

var signalBiz = SignalBiz{}

// Mock 根据给定的id模拟信号延迟警告业务的规则执行。
// 该方法主要用于测试或模拟特定规则的执行，以验证规则的有效性和触发条件。
//
// 参数:
//
//	id - 规则的唯一标识符。
//
// 返回值:
//
//	bool - 表示模拟执行的结果，通常指示规则是否被触发。`
func (biz *SignalDelayWaringBiz) Mock(id int) bool {
	rule, mm := biz.GenParam(id)
	glob.GLog.Sugar().Infof("模拟参数 %+v", mm)

	vm := goja.New()
	_, err := vm.RunString(rule.Script)
	if err != nil {
		fmt.Println("JS代码有问题！")
	}
	var fn func(string2 map[string][]v) bool
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		fmt.Println("Js函数映射到 Go 函数失败！")
		panic(err)
	}
	a := fn(mm)
	return a

}

// GenParam 根据给定的id生成信号延迟警告参数。
// 该方法检索与id匹配的信号延迟警告规则，并为该规则的每个参数生成模拟数据。
//
// 参数:
//
//   - id - 信号延迟警告规则的唯一标识。
//
// 返回值:
//
//   - models.SignalDelayWaring - 与给定id匹配的信号延迟警告规则。
//   - map[string][]v - 参数名到模拟数据值的映射，其中v的类型是一个包含时间和值的结构体。
func (biz *SignalDelayWaringBiz) GenParam(id int) (models.SignalDelayWaring, map[string][]v) {
	var dt models.SignalDelayWaring
	db := glob.GDb
	db.First(&dt, id)

	var param []models.SignalDelayWaringParam
	db.Model(models.SignalDelayWaringParam{}).Where("signal_delay_waring_id = ?", id).Find(&param)

	mm := make(map[string][]v)

	for _, waringParam := range param {
		signal, err := signalBiz.FindByIdForSignal(waringParam.SignalId)
		if err != nil {
			glob.GLog.Sugar().Errorf("查询异常 %+v", err)
		}
		vm := []v{}

		for range signal.CacheSize {
			a := v{
				Time:  time.Now().Unix(),
				Value: 10,
			}
			vm = append(vm, a)
		}
		mm[waringParam.Name] = vm
	}
	return dt, mm
}

func (b SignalDelayWaringBiz) InitMongoCollection(m *models.SignalDelayWaring) {
	name := ut.CalcCollectionName(glob.GConfig.MongoConfig.ScriptWaringCollection, m.ID)
	ut.CheckCollectionAndCreate(glob.GConfig.MongoConfig.ScriptWaringCollection, name)
}

type v struct {
	Time  int64   `json:"Time"`
	Value float64 `json:"Value"`
}
