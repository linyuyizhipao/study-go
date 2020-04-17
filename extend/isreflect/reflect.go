package isreflect

import (
	"fmt"
	"reflect"
)

//断言 获取结构体的每一个key和值
//reflect 得到的tt 可以使用reflect的方法与reflect里面的常量比较，从而判断得到类型
//核心思想就是，reflect一旦使用，就使用到底
func ReadStruct(s interface{}) (m map[string]interface{}){
	m = make(map[string]interface{})
	tt :=reflect.TypeOf(s) //key操作对象
	va :=reflect.ValueOf(s)  //value 操作对象

	switch tt.Kind() {
	case reflect.Struct:
		processStruct(va,tt)
	case reflect.Int:
		processInt(va,tt)
	}
	return
}

func processStruct(va reflect.Value,tt reflect.Type){
	for i:=0;i<va.NumField();i++{
		ke :=tt.Field(i).Name
		fmt.Println("我是结构体key:%s,我是结构体值:%s",ke,va.Interface())
	}
}

func processInt(va reflect.Value,tt reflect.Type){
	ke :=tt.Name
	fmt.Println("我是int的key:%s,我是结int:%d",ke,va.Interface())
}
