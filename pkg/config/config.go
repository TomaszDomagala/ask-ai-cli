package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

func FillWithConfig(config *viper.Viper, input interface{}) error {
	return config.Unmarshal(input, func(config *mapstructure.DecoderConfig) {
		config.TagName = "name"
	})
}

func BindProviderFlags(conf *viper.Viper, flags *pflag.FlagSet, provider string) error {
	var firstErr error
	prefix := provider + "-"

	flags.VisitAll(func(flag *pflag.Flag) {
		if strings.HasPrefix(flag.Name, prefix) {
			if err := conf.BindPFlag(strings.TrimPrefix(flag.Name, prefix), flag); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	})
	return firstErr
}

// FlagsFromStruct adds flags to the provided flags set from the struct tags.
// It reads the tags: name, shorthand, value, usage.
// prefix is added to the flag name with a dash separator.
func FlagsFromStruct(flags *pflag.FlagSet, input interface{}, prefix string) {
	tags := structTags(input, []string{"name", "shorthand", "value", "usage"})
	types := structFieldType(input)

	t := reflect.TypeOf(input)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldType := types[fieldName]
		fieldTags := tags[fieldName]
		name := fieldTags["name"]
		shorthand := fieldTags["shorthand"]
		value := fieldTags["value"]
		usage := fieldTags["usage"]

		if prefix != "" {
			name = prefix + "-" + name
		}

		switch fieldType.Kind() {
		case reflect.String:
			flags.StringP(name, shorthand, cast.ToString(value), usage)
		case reflect.Int:
			flags.IntP(name, shorthand, cast.ToInt(value), usage)
		case reflect.Float64:
			flags.Float64P(name, shorthand, cast.ToFloat64(value), usage)
		case reflect.Struct:
			FlagsFromStruct(flags, reflect.New(fieldType).Elem().Interface(), prefix)
		default:
			panic(fmt.Sprintf("unsupported type %s, field %s", fieldType.Kind(), fieldName))
		}
	}
}

// structTags returns tags values per field name of a struct
func structTags(input interface{}, tags []string) map[string]map[string]string {
	t := reflect.TypeOf(input)
	sTags := make(map[string]map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldTags := make(map[string]string)
		for _, tag := range tags {
			fieldTags[tag] = field.Tag.Get(tag)
		}
		sTags[field.Name] = fieldTags
	}
	return sTags
}

// structFieldType returns field type per field name of a struct
func structFieldType(input interface{}) map[string]reflect.Type {
	t := reflect.TypeOf(input)
	sTypes := make(map[string]reflect.Type)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sTypes[field.Name] = field.Type
	}
	return sTypes
}
