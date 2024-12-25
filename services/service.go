package services

import (
	"github.com/leor-w/injector"
	"reflect"
)

func InitService(appScope *injector.Scope) {
	appScope.Provide(new(LoginService))
	appScope.Provide(new(AdminService))
	appScope.Provide(new(UploadService))
	appScope.Provide(new(SmsService))
	appScope.Provide(new(UserService))
	appScope.Provide(new(RoleService))
}

// GetFieldValues 使用反射获取结构体切片中某个字段的全部值
func GetFieldValues(slice interface{}, fieldName string) []interface{} {
	// 获取切片的反射值
	sliceValue := reflect.ValueOf(slice)

	// 确保传入的是切片
	if sliceValue.Kind() != reflect.Slice {
		panic("getFieldValues: not a slice")
	}

	// 获取字段的反射类型
	_, exist := sliceValue.Type().Elem().FieldByName(fieldName)

	// 如果字段不存在，则返回空切片
	if !exist {
		return []interface{}{}
	}

	// 定义一个切片来存储字段值
	fieldValues := make([]interface{}, sliceValue.Len())

	// 遍历切片，获取每个结构体的字段值
	for i := 0; i < sliceValue.Len(); i++ {
		fieldValue := sliceValue.Index(i).FieldByName(fieldName).Interface()
		fieldValues[i] = fieldValue
	}

	return fieldValues
}
