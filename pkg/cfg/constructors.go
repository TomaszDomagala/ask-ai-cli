package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// newWithFlagSetter creates a new config Value with the provided flag setter.
// Flag setter is prepended to the options to ensure that it is applied first.
func newWithFlagSetter[T any](config *viper.Viper, key string, getter Getter[T], flagSetter FlagSetter[T], options ...Option[T]) Value[T] {
	return New(config, key, getter, append([]Option[T]{WithFlagSetter(flagSetter)}, options...)...)
}

// String creates a new config Value of type string.
func String(config *viper.Viper, key string, options ...Option[string]) Value[string] {
	return newWithFlagSetter(config, key, (*viper.Viper).GetString, (*pflag.FlagSet).StringP, options...)
}

// Int creates a new config Value of type int.
func Int(config *viper.Viper, key string, options ...Option[int]) Value[int] {
	return newWithFlagSetter(config, key, (*viper.Viper).GetInt, (*pflag.FlagSet).IntP, options...)
}
