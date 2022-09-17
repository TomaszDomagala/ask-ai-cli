// Package config is a wrapper around viper that provides
// a consistent way to access configuration values.
package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Getter is a function that returns a value from the config.
type Getter[T any] func(*viper.Viper, string) T

// FlagSetter is a function that creates a flag for the provided flag set.
type FlagSetter[T any] func(f *pflag.FlagSet, name, shorthand string, value T, usage string) *T

// Value is a wrapper around a config value.
type Value[T any] struct {
	// config is the viper instance that holds the configuration.
	config *viper.Viper
	// key is the key of the value in the config.
	key string
	// getter is a function that returns the value associated with the key as a type T.
	getter Getter[T]
	// flagSetter is a function that creates a flag for the provided flag set.
	flagSetter FlagSetter[T]
	// flag is the optional flag associated with the value.
	flag *pflag.Flag
	// flagValue is the optional flag value associated with the value.
	flagValue *T

	// postAttach is a list of functions that are executed after Attach method is called.
	postAttach []func() error
}

type Option[T any] func(*Value[T])

// WithFlagSetter sets the flag setter for the value.
func WithFlagSetter[T any](flagSetter FlagSetter[T]) Option[T] {
	return func(v *Value[T]) {
		v.flagSetter = flagSetter
	}
}

// WithFlag creates a flag for the value.
// If used with New, remember to set use WithFlagSetter before.
// The flag will be bound to the config when Attach is called.
func WithFlag[T any](f *pflag.FlagSet, name, shorthand string, value T, usage string) Option[T] {
	return func(v *Value[T]) {
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

// New creates a new config Value.
func New[T any](key string, getter Getter[T], options ...Option[T]) Value[T] {
	value := Value[T]{
		key:    key,
		getter: getter,
	}
	for _, option := range options {
		option(&value)
	}
	return value
}

func (v Value[T]) Get() T {
	return v.getter(v.config, v.key)
}

// Attach attaches config to the Value and executes post attach functions,
// such as binding the flag to the config.
func (v Value[any]) Attach(config *viper.Viper) error {
	v.config = config

	for _, postAttach := range v.postAttach {
		if err := postAttach(); err != nil {
			return err
		}
	}
	return nil
}
