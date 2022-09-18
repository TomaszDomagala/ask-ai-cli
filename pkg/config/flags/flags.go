// Package flags provides a wrapper around pflag that provides
// a consistent way to access configuration values.
package flags

import (
	"github.com/spf13/pflag"
)

type Flag[T any] struct {
	// name is the name of the flag.
	name string
	// flags is the flag set that the flag is attached to.
	flags *pflag.FlagSet
	// flag is the pflag.Flag associated with the Flag.
	flag *pflag.Flag
	// flagValue is the value associated with the Flag.
	flagValue *T
}

// Get returns the value of the flag.
func (f Flag[T]) Get() T {
	return *f.flagValue
}

// Changed returns true if the flag was explicitly set during Parse() and false
// otherwise
func (f Flag[T]) Changed() bool {
	return f.flags.Changed(f.name)
}

// flagSetter is a function from pflag.FlagSet that creates the flag.
// Example: (*pflag.FlagSet).StringP, (*pflag.FlagSet).IntP, etc.
type flagSetter[T any] func(f *pflag.FlagSet, name, shorthand string, value T, usage string) *T

// newFlagP creates a new flag with the provided name, shorthand, value and usage.
func newFlagP[T any](flags *pflag.FlagSet, setter flagSetter[T], name, shorthand string, value T, usage string) Flag[T] {
	flagValue := setter(flags, name, shorthand, value, usage)
	flag := flags.Lookup(name)

	return Flag[T]{
		name:      name,
		flag:      flag,
		flagValue: flagValue,
		flags:     flags,
	}
}

// newFlag creates a new flag with the provided name, value and usage.
func newFlag[T any](flags *pflag.FlagSet, setter flagSetter[T], name string, value T, usage string) Flag[T] {
	return newFlagP(flags, setter, name, "", value, usage)
}

// String defines a string flag with specified name, default value, and usage string.
func String(flags *pflag.FlagSet, name string, value string, usage string) Flag[string] {
	return newFlag(flags, (*pflag.FlagSet).StringP, name, value, usage)
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func StringP(flags *pflag.FlagSet, name, shorthand string, value string, usage string) Flag[string] {
	return newFlagP(flags, (*pflag.FlagSet).StringP, name, shorthand, value, usage)
}

// Bool defines a bool flag with specified name, default value, and usage string.
func Bool(flags *pflag.FlagSet, name string, value bool, usage string) Flag[bool] {
	return newFlag(flags, (*pflag.FlagSet).BoolP, name, value, usage)
}

// BoolP is like Bool, but accepts a shorthand letter that can be used after a single dash.
func BoolP(flags *pflag.FlagSet, name, shorthand string, value bool, usage string) Flag[bool] {
	return newFlagP(flags, (*pflag.FlagSet).BoolP, name, shorthand, value, usage)
}
