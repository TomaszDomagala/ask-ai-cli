package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// attachable is an interface that can be used to attach a config to a configValue.
type attachable interface {
	Attach(config *viper.Viper) error
}

// Attach attaches config to the configValue and executes post attach functions,
// such as binding the flag to the config.
func (v *configValue[T]) Attach(config *viper.Viper) error {
	v.config = config

	for _, postAttach := range v.postAttach {
		if err := postAttach(); err != nil {
			return err
		}
	}
	return nil
}

// Attach attaches the config to all fields of the provided struct.
// The input must be a pointer to a struct.
// To attach a config to a field, the field must be of attachable interface
// and be exported.
func Attach(config *viper.Viper, input interface{}) error {
	return Traverse(input, func(value AnyValue) error {
		if att, ok := value.(attachable); ok {
			err := att.Attach(config)
			if err != nil {
				return fmt.Errorf("failed to attach config to %s: %w", value.Key(), err)
			}
		}
		return nil
	})
}
