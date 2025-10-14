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
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"igp/servlet"
)

type ScriptBiz struct{}

func (biz *ScriptBiz) CheckScript(param string, script string) *[]servlet.DataRowList {

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		zap.S().Errorf("JS代码有问题！ %+v", err)
		return nil
	}
	var fn func(string2 string) *[]servlet.DataRowList
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		zap.S().Errorf("Js函数映射到 Go 函数失败！ %+v", err)
		return nil
	}
	a := fn(param)
	return a

}
