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
	"github.com/dop251/goja"
	"igp/servlet"
	"testing"
)

func TestFibScript(t *testing.T) {
	// 您的JavaScript代码
	const script = `
        function main(nc) {
            var dataRows = [
                { "Name": "Temperature", "Value": "23" },
                { "Name": "Humidity", "Value": "30" }
            ];
            var result = {
                "Time": 1652499200000,
                "DataRows": dataRows,
                "Nc": nc // 确保结果对象中包含nc参数
            };
            return result;
        }
    `

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		fmt.Println("JS代码有问题！")
		return
	}
	var fn func(string2 string) *servlet.DataRowList
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		fmt.Println("Js函数映射到 Go 函数失败！")
		return
	}
	a := fn("aaa")
	fmt.Println("执行：", a)

}
