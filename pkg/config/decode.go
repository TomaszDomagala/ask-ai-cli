package config

import (
	"fmt"
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
	m := make(map[string]any)
	err := Traverse(input, func(value AnyValue) error {
		m[value.Key()] = value.GetAny()
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to convert input to map: %w", err)
	}

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
