package config

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// Decode takes a struct as input end transforms the Value fields
// into the output struct fields. It uses mapstructure to decode the
// map into the output struct. It is important to remember that the input
// field's key should be the same as the `config` tag of the output field.
// Example:
//
//	type config struct {
//		ApiKey Value[string]
//	}
//	type concreteConfig struct {
//		ApiKey string `config:"api_key"`
//	}
//	var cfg config = config{
//		ApiKey: config.String("api_key"),
//	}
//	var concreteCfg concreteConfig
//	err := Flatten(cfg, &concreteCfg)
func Decode(input interface{}, output interface{}) error {
	m, err := toMap(input)
	if err != nil {
		return fmt.Errorf("failed to convert input to map: %w", err)
	}

	fmt.Println(m)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  output,
		TagName: "config",
		Squash:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}
	if err = decoder.Decode(m); err != nil {
		return fmt.Errorf("failed to decode output: %w", err)
	}
	return nil
}

// anyValue is a weakly typed Value interface.
type anyValue interface {
	Key() string
	getAny() any
}

// getAny returns the value associated with the key as a type any.
func (v *configValue[T]) getAny() any {
	return v.getter(v.config, v.key)
}

// toMap converts the config to a map.
func toMap(input any) (map[string]any, error) {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct, got %T", input)
	}
	t := v.Type()
	m := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		if field.Anonymous {
			embedded, err := toMap(v.Field(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to convert embedded field %q: %w", field.Name, err)
			}
			for key, value := range embedded {
				m[key] = value
			}
		}
		if value, ok := v.Field(i).Interface().(anyValue); ok {
			m[value.Key()] = value.getAny()
		}
	}
	return m, nil
}
