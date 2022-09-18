package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

// attachable is an interface that can be used to attach a config to a configValue.
type attachable interface {
	Attach(config *viper.Viper) error
}

// Attach attaches the config to all fields of the provided struct.
// The input must be a pointer to a struct.
// To attach a config to a field, the field must be of attachable interface
// and be exported.
func Attach(config *viper.Viper, input interface{}) error {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct or a pointer to a struct, got %T", input)
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}
		if field.Anonymous {
			if err := Attach(config, v.Field(i).Addr().Interface()); err != nil {
				return fmt.Errorf("failed to attach embedded field %q: %w", field.Name, err)
			}
			continue
		}

		// try to get the attachable interface, if it exists, use it
		// otherwise, skip the field
		if att, ok := v.Field(i).Interface().(attachable); ok {
			if err := att.Attach(config); err != nil {
				return fmt.Errorf("failed to attach field %q: %w", field.Name, err)
			}
		}
	}

	return nil
}
