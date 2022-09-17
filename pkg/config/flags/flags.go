package flags

import (
	"github.com/spf13/pflag"
)

// FlagSetter is a function that creates a flag for the provided flag set.
type FlagSetter[T any] func(f *pflag.FlagSet, name, shorthand string, value T, usage string) *T

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

func (f *Flag[T]) Get() T {
	return *f.flagValue
}

func New[T any](flags *pflag.FlagSet, setter FlagSetter[T], name, shorthand string, value T, usage string) *Flag[T] {
	flagValue := setter(flags, name, shorthand, value, usage)
	flag := flags.Lookup(name)

	return &Flag[T]{
		name:      name,
		flag:      flag,
		flagValue: flagValue,
		flags:     flags,
	}
}

func String(flags *pflag.FlagSet, name, shorthand string, value string, usage string) *Flag[string] {
	return New(flags, (*pflag.FlagSet).StringP, name, shorthand, value, usage)
}

func Int(flags *pflag.FlagSet, name, shorthand string, value int, usage string) *Flag[int] {
	return New(flags, (*pflag.FlagSet).IntP, name, shorthand, value, usage)
}

func Bool(flags *pflag.FlagSet, name, shorthand string, value bool, usage string) *Flag[bool] {
	return New(flags, (*pflag.FlagSet).BoolP, name, shorthand, value, usage)
}
