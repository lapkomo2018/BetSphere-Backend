package validate

import (
	"fmt"
	"reflect"
)

func StructPointersNotNil(cfg interface{}) error {
	val := reflect.ValueOf(cfg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", val.Kind())
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			return fmt.Errorf("field %q is nil", fieldType.Name)
		}
	}

	return nil
}
