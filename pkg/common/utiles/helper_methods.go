package utiles

import (
	"reflect"
)

func ContainsStructFieldValue(slice interface{}, fieldName string, fieldValueToCheck interface{}) bool {
	rangeOnMe := reflect.ValueOf(slice)
	for i := 0; i < rangeOnMe.Len(); i++ {
		s := rangeOnMe.Index(i)
		f := s.FieldByName(fieldName)
		if f.IsValid() {
			if f.Interface() == fieldValueToCheck {
				return true
			}
		}
	}
	return false
}
