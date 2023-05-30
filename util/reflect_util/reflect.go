package reflect_util

import "reflect"

func GetFields(t reflect.Type) []reflect.StructField {
	// 获取给定类型的字段列表
	// 获取字段数量
	num := t.NumField() 

	// 创建存储字段的切片
	fields := make([]reflect.StructField, 0, num) 
	for i := 0; i < num; i++ {
		// 将字段添加到切片中
		fields = append(fields, t.Field(i)) 
	}

	return fields
}

func GetFieldValues(v reflect.Value) []reflect.Value {
	// 获取给定值的字段值列表
	num := v.NumField() 

	// 获取字段数量
	fieldValues := make([]reflect.Value, 0, num) 
	// 创建存储字段值的切片
	for i := 0; i < num; i++ {
		// 将字段值添加到切片中
		fieldValues = append(fieldValues, v.Field(i)) 
	}

	return fieldValues
}