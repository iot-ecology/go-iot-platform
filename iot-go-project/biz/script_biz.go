package biz

import (
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"igp/servlet"
)

type ScirptBiz struct{}

func (biz *ScirptBiz) CheckScript(param string, script string) *[]servlet.DataRowList {

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
