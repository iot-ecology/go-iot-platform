package biz

import (
	"fmt"
	"github.com/dop251/goja"
	"igp/servlet"
)

type ScirptBiz struct{}

func (biz *ScirptBiz) CheckScript(param string, script string) *[]servlet.DataRowList {

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		fmt.Println("JS代码有问题！")
		return nil
	}
	var fn func(string2 string) *[]servlet.DataRowList
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		fmt.Println("Js函数映射到 Go 函数失败！")
		return nil
	}
	a := fn(param)
	return a

}
