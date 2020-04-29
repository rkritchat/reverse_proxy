package common

import "reflect"

func ConvertToArray(slice interface{}) []interface{}{
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	r := make([]interface{}, s.Len())
	for i:=0; i<s.Len(); i++ {
		r[i] = s.Index(i).Interface()
	}
	return r
}
