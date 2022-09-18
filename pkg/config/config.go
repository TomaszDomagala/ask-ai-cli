// Package config is a wrapper around viper that provides
// a consistent way to access configuration values.
package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Value is an interface that can be used to get values from a config.
type Value[T any] interface {
	Key() string
	Get() T
}

// viperGetter is a function that returns a value from the viper config.
type viperGetter[T any] func(*viper.Viper, string) T

// pflagSetter is a function that creates a flag for the provided flag set.
type pflagSetter[T any] func(f *pflag.FlagSet, name, shorthand string, value T, usage string) *T

// configValue is a wrapper around a viper config value.
type configValue[T any] struct {
	// config is the viper instance that holds the configuration.
	config *viper.Viper
	// key is the key of the value in the config.
	key string
	// getter is a function that returns the value associated with the key as a type T.
	getter viperGetter[T]
	// flagSetter is a function that creates a flag for the provided flag set.
	flagSetter pflagSetter[T]
	// flag is the optional flag associated with the value.
	flag *pflag.Flag
	// flagValue is the optional flag value associated with the value.
	flagValue *T

	// postAttach is a list of functions that are executed after Attach method is called.
	postAttach []func() error
}

func (v *configValue[T]) Get() T {
	return v.getter(v.config, v.key)
}

func (v *configValue[T]) Key() string {
	return v.key
}

type Option[T any] func(*configValue[T])

// WithFlagP is like WithFlag, but accepts a shorthand letter that can be used after a single dash.
func WithFlagP[T any](f *pflag.FlagSet, name, shorthand string, value T, usage string) Option[T] {
	return func(v *configValue[T]) {
		v.flagValue = v.flagSetter(f, name, shorthand, value, usage)
		v.flag = f.Lookup(name)

		v.postAttach = append(v.postAttach, func() error {
			if err := v.config.BindPFlag(v.key, v.flag); err != nil {
				return fmt.Errorf("failed to bind flag %q to config key %q: %w", v.flag.Name, v.key, err)
			}
			return nil
		})
	}
}

// WithFlag defines a flag with specified name, default value, and usage string.
func WithFlag[T any](f *pflag.FlagSet, name string, value T, usage string) Option[T] {
	return WithFlagP(f, name, "", value, usage)
}

// Attach attaches config to the configValue and executes post attach functions,
// such as binding the flag to the config.
func (v *configValue[any]) Attach(config *viper.Viper) error {
	v.config = config

	for _, postAttach := range v.postAttach {
		if err := postAttach(); err != nil {
			return err
		}
	}
	return nil
}

// newValue creates a new config configValue.
func newValue[T any](key string, getter viperGetter[T], flagSetter pflagSetter[T], options ...Option[T]) *configValue[T] {
	value := configValue[T]{
		key:        key,
		getter:     getter,
		flagSetter: flagSetter,
	}
	for _, option := range options {
		option(&value)
	}
	return &value
}

// String creates a new config configValue of type string.
func String(key string, options ...Option[string]) Value[string] {
	return newValue(key, (*viper.Viper).GetString, (*pflag.FlagSet).StringP, options...)
}

// Int creates a new config configValue of type int.
func Int(key string, options ...Option[int]) Value[int] {
	return newValue(key, (*viper.Viper).GetInt, (*pflag.FlagSet).IntP, options...)
}

// Float64 creates a new config configValue of type float64.
func Float64(key string, options ...Option[float64]) Value[float64] {
	return newValue(key, (*viper.Viper).GetFloat64, (*pflag.FlagSet).Float64P, options...)
}
