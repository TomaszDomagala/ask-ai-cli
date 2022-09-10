package config

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestStructTags(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		tags   []string
		output map[string]map[string]string
	}{
		{
			name: "single field",
			input: struct {
				Foo string `json:"foo"`
			}{},
			tags: []string{"json"},
			output: map[string]map[string]string{
				"Foo": {
					"json": "foo",
				},
			},
		},
		{
			name: "multiple fields",
			input: struct {
				Foo string `json:"foo"`
				Bar string `json:"bar" yarn:"bar"`
			}{},
			tags: []string{"json", "yarn"},
			output: map[string]map[string]string{
				"Foo": {
					"json": "foo",
					"yarn": "",
				},
				"Bar": {
					"json": "bar",
					"yarn": "bar",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := structTags(test.input, test.tags)

			if !reflect.DeepEqual(got, test.output) {
				t.Errorf("got %v, want %v", got, test.output)
			}
		})
	}
}

func TestFlagsFromStruct(t *testing.T) {
	type Flag struct {
		name      string
		shorthand string
		value     string
		usage     string
	}

	tests := []struct {
		name   string
		input  interface{}
		prefix string
		flags  []Flag
	}{
		{
			name: "single field",
			input: struct {
				Foo string `name:"foo" shorthand:"f" value:"bar" usage:"foo bar"`
			}{},
			flags: []Flag{
				{
					name:      "foo",
					shorthand: "f",
					value:     "bar",
					usage:     "foo bar",
				},
			},
		},
		{
			name: "multiple fields",
			input: struct {
				Foo string `name:"foo" shorthand:"f" value:"bar" usage:"foo bar"`
				Bar string `name:"bar" shorthand:"b" value:"baz" usage:"bar baz"`
			}{},
			flags: []Flag{
				{
					name:      "foo",
					shorthand: "f",
					value:     "bar",
					usage:     "foo bar",
				},
				{
					name:      "bar",
					shorthand: "b",
					value:     "baz",
					usage:     "bar baz",
				},
			},
		},
		{
			name: "with prefix",
			input: struct {
				Foo string `name:"foo" shorthand:"f" value:"bar" usage:"foo bar"`
			}{},
			prefix: "prefix",
			flags: []Flag{
				{
					name:      "prefix-foo",
					shorthand: "f",
					value:     "bar",
					usage:     "foo bar",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			FlagsFromStruct(flags, test.input, test.prefix)
			for _, flag := range test.flags {
				f := flags.Lookup(flag.name)
				require.NotNil(t, f)

				assert.Equal(t, flag.name, f.Name)
				assert.Equal(t, flag.shorthand, f.Shorthand)
				assert.Equal(t, flag.value, f.Value.String())
				assert.Equal(t, flag.usage, f.Usage)
			}
		})
	}
}
