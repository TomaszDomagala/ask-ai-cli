package config

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrInvalidInput is returned when the input is not a pointer to a struct or a struct.
	ErrInvalidInput = errors.New("input must be a struct or a pointer to a struct")
)

// AnyValue is a weakly typed Value interface.
type AnyValue interface {
	Key() string
	IsSet() bool

	GetAny() any
	SetAny(any)
}

// GetAny is a helper function that returns the value associated with the key as a type any.
func (v *configValue[T]) GetAny() any {
	return v.Get()
}

// SetAny is a helper function that sets the value associated with the key as a type any.
// Note that it does not check if the provided value's is correct.
func (v *configValue[T]) SetAny(value any) {
	v.config.Set(v.key, value)
}

// Traverse traverses the input struct recursively and calls
// the provided function for each AnyValue field.
// It stops if the function returns an error.
// Input must be a pointer to a struct or a struct.
func Traverse(input interface{}, f func(AnyValue) error) error {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("%w, got %T", ErrInvalidInput, input)
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldValue := v.Field(i)
		if field.Anonymous {
			if err := Traverse(fieldValue.Interface(), f); err != nil {
				return fmt.Errorf("failed to traverse embedded field %q: %w", field.Name, err)
			}
		}
		if value, ok := fieldValue.Interface().(AnyValue); ok {
			if err := f(value); err != nil {
				return fmt.Errorf("failed to traverse field %q: %w", field.Name, err)
			}
		}
	}
	return nil
}
