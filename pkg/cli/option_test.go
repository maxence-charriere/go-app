package cli

import (
	"flag"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type nestedStruct struct {
	Float float64
}

func TestOptionParserParse(t *testing.T) {
	tests := []struct {
		scenario        string
		args            []string
		env             map[string]string
		options         interface{}
		expectedOptions interface{}
		err             bool
		parseFlagsErr   bool
	}{
		{
			scenario: "parsing a non-pointer option returns an error",
			options:  struct{}{},
			err:      true,
		},
		{
			scenario: "parsing a non-struct pointer option returns an error",
			options:  new(int),
			err:      true,
		},
		{
			scenario: "parsing a nil struct pointer returns an error",
			options:  (*struct{})(nil),
			err:      true,
		},
		{
			scenario: "parsing a nil scalar pointer returns an error",
			options:  (*int)(nil),
			err:      true,
		},
		{
			scenario:        "parsing empty options succeed",
			options:         &struct{}{},
			expectedOptions: struct{}{},
		},
		{
			scenario: "parsing options succeed",
			options: &struct {
				Int    int
				String string
			}{
				Int:    42,
				String: "foo",
			},
			expectedOptions: struct {
				Int    int
				String string
			}{
				Int:    42,
				String: "foo",
			},
		},
		{
			scenario: "parsing options from env succeed",
			env: map[string]string{
				"INT":    "21",
				"STRING": "bar",
			},
			options: &struct {
				Int    int
				String string
			}{
				Int:    42,
				String: "foo",
			},
			expectedOptions: struct {
				Int    int
				String string
			}{
				Int:    21,
				String: "bar",
			},
		},
		{
			scenario: "parsing options from env with tagged name succeed",
			env: map[string]string{
				"INT":         "21",
				"TEST_STRING": "bar",
			},
			options: &struct {
				Int    int
				String string `env:"TEST_STRING"`
			}{
				Int:    42,
				String: "foo",
			},
			expectedOptions: struct {
				Int    int
				String string `env:"TEST_STRING"`
			}{
				Int:    21,
				String: "bar",
			},
		},
		{
			scenario: "ignore env",
			env: map[string]string{
				"-":      "21",
				"STRING": "bar",
			},
			options: &struct {
				Int    int `env:"-"`
				String string
			}{
				Int:    42,
				String: "foo",
			},
			expectedOptions: struct {
				Int    int `env:"-"`
				String string
			}{
				Int:    42,
				String: "bar",
			},
		},
		{
			scenario: "parsing options from args succeed",
			args: []string{
				"-int", "84",
				"-string", "boo",
			},
			options: &struct {
				Int    int
				String string
			}{
				Int:    42,
				String: "foo",
			},
			expectedOptions: struct {
				Int    int
				String string
			}{
				Int:    84,
				String: "boo",
			},
		},
		{
			scenario: "parsing options from args with tagged name succeed",
			args: []string{
				"-i", "21",
				"-string", "boo",
			},
			options: &struct {
				Int    int `cli:"i"`
				String string
			}{
				Int:    42,
				String: "foo",
			},
			expectedOptions: struct {
				Int    int `cli:"i"`
				String string
			}{
				Int:    21,
				String: "boo",
			},
		},
		{
			scenario: "nonexported args are ignored",
			args: []string{
				"-int", "84",
			},
			options: &struct {
				Int int
				str string
			}{
				Int: 42,
			},
			expectedOptions: struct {
				Int int
				str string
			}{
				Int: 84,
			},
		},
		{
			scenario: "args take priority over env variables",
			env: map[string]string{
				"INT":    "21",
				"STRING": "bar",
			},
			args: []string{
				"-int", "84",
				"-string", "boo",
			},
			options: &struct {
				Int    int
				String string
			}{
				Int:    42,
				String: "foo",
			},
			expectedOptions: struct {
				Int    int
				String string
			}{
				Int:    84,
				String: "boo",
			},
		},
		{
			scenario: "parsing options with bool values succeed",
			args: []string{
				"-f",
			},
			options: &struct {
				Force   bool `cli:"f"`
				Verbose bool `cli:"v"`
			}{},
			expectedOptions: struct {
				Force   bool `cli:"f"`
				Verbose bool `cli:"v"`
			}{
				Force: true,
			},
		},
		{
			scenario: "parsing options with nested struct succeed",
			args: []string{
				"-int", "21",
				"-struct", `{"Float":49.3}`,
			},
			options: &struct {
				Int    int
				Struct nestedStruct
			}{},
			expectedOptions: struct {
				Int    int
				Struct nestedStruct
			}{
				Int: 21,
				Struct: nestedStruct{
					Float: 49.3,
				},
			},
		},
		{
			scenario: "parsing options with defined nested struct fields succeed",
			args: []string{
				"-int", "21",
				"-struct.float", "49.3",
			},
			options: &struct {
				Int    int
				Struct nestedStruct
			}{},
			expectedOptions: struct {
				Int    int
				Struct nestedStruct
			}{
				Int: 21,
				Struct: nestedStruct{
					Float: 49.3,
				},
			},
		},
		// {
		// 	scenario: "parsing options with invalid nested json returns an error",
		// 	args: []string{
		// 		"-int", "21",
		// 		"-struct.float", "{}",
		// 	},
		// 	options: &struct {
		// 		Int    int
		// 		Struct nestedStruct
		// 	}{},
		// 	parseFlagsErr: true,
		// },
		{
			scenario: "parsing options with a slice field succeed",
			args: []string{
				"-slice", `["foo","bar"]`,
			},
			options: &struct {
				Slice []string
			}{},
			expectedOptions: struct {
				Slice []string
			}{
				Slice: []string{"foo", "bar"},
			},
		},
		{
			scenario: "parsing options with duration fields succeed",
			args: []string{
				"-duration", "42",
			},
			options: &struct {
				Duration time.Duration
			}{},
			expectedOptions: struct {
				Duration time.Duration
			}{
				Duration: 42,
			},
		},
		{
			scenario: "parsing options with duration fields and litteral format succeed",
			args: []string{
				"-duration", "42s",
			},
			options: &struct {
				Duration time.Duration
			}{},
			expectedOptions: struct {
				Duration time.Duration
			}{
				Duration: 42000000000,
			},
		},
		{
			scenario: "parsing options with invalid duration fields return an error",
			args: []string{
				"-duration", "4klm41238rub+_8u498qurfvn",
			},
			options: &struct {
				Duration time.Duration
			}{},
			parseFlagsErr: true,
		},
		{
			scenario: "parsing options with time fields succeed",
			args: []string{
				"-time", "1986-02-14T00:00:00Z",
			},
			options: &struct {
				Time time.Time
			}{},
			expectedOptions: struct {
				Time time.Time
			}{
				Time: time.Date(1986, 2, 14, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			flags := flag.NewFlagSet("test", flag.ContinueOnError)
			flags.SetOutput(io.Discard)

			p := optionParser{
				flags: flags,
			}

			for k, v := range test.env {
				t.Setenv(k, v)
			}

			_, err := p.parse(test.options)
			if test.err {
				require.Error(t, err)
				t.Log("error:", err)
				return
			}
			require.NoError(t, err)

			if test.parseFlagsErr {
				err = p.flags.Parse(test.args)
				require.Error(t, err)
				return
			}

			err = p.flags.Parse(test.args)
			require.NoError(t, err)

			options := reflect.ValueOf(test.options).Elem().Interface()
			require.Equal(t, test.expectedOptions, options)
		})
	}
}

func TestOptionParserNil(t *testing.T) {
	p := optionParser{flags: flag.NewFlagSet("test", flag.ContinueOnError)}
	opts, err := p.parse(nil)
	require.NoError(t, err)
	require.Empty(t, opts)
}

func TestOptionParserInvalidEnvironment(t *testing.T) {
	const envKey = "GOAPP_TEST_OPTION"
	const secret = "private-config-value"
	tests := []struct {
		name       string
		options    any
		value      string
		optionName string
	}{
		{
			name: "integer",
			options: &struct {
				Value int `env:"GOAPP_TEST_OPTION"`
			}{Value: 42},
			value:      `"` + secret + `"`,
			optionName: "value",
		},
		{
			name: "boolean",
			options: &struct {
				Value bool `env:"GOAPP_TEST_OPTION"`
			}{},
			value:      `"` + secret + `"`,
			optionName: "value",
		},
		{
			name: "duration",
			options: &struct {
				Value time.Duration `env:"GOAPP_TEST_OPTION"`
			}{},
			value:      secret,
			optionName: "value",
		},
		{
			name: "time",
			options: &struct {
				Value time.Time `env:"GOAPP_TEST_OPTION"`
			}{},
			value:      secret,
			optionName: "value",
		},
		{
			name: "time pointer",
			options: &struct {
				Value *time.Time `env:"GOAPP_TEST_OPTION"`
			}{},
			value:      secret,
			optionName: "value",
		},
		{
			name: "slice",
			options: &struct {
				Value []int `env:"GOAPP_TEST_OPTION"`
			}{},
			value:      `["` + secret + `"]`,
			optionName: "value",
		},
		{
			name: "map",
			options: &struct {
				Value map[string]int `env:"GOAPP_TEST_OPTION"`
			}{},
			value:      `{"key":"` + secret + `"}`,
			optionName: "value",
		},
		{
			name: "struct",
			options: &struct {
				Value nestedStruct `env:"GOAPP_TEST_OPTION"`
			}{},
			value:      `{"Float":"` + secret + `"}`,
			optionName: "value",
		},
		{
			name: "nested field",
			options: &struct {
				Nested struct {
					Count int `env:"GOAPP_TEST_OPTION"`
				} `env:"-"`
			}{},
			value:      `"` + secret + `"`,
			optionName: "nested.count",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(envKey, tt.value)
			p := optionParser{flags: flag.NewFlagSet("test", flag.ContinueOnError)}
			_, err := p.parse(tt.options)
			require.Error(t, err)
			require.Contains(t, err.Error(), envKey)
			require.Contains(t, err.Error(), tt.optionName)
			require.NotContains(t, err.Error(), secret)
		})
	}
}

func TestNormalizeOptionName(t *testing.T) {
	tests := []struct {
		scenario     string
		baseName     string
		expectedName string
	}{
		{
			scenario:     "name with camel case",
			baseName:     "HelloWorld",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with number",
			baseName:     "HelloWorld42",
			expectedName: "hello-world42",
		},
		{
			scenario:     "name with dash",
			baseName:     "hello-world",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with dash and upper case",
			baseName:     "hello-World",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with dash at the end",
			baseName:     "hello-World-",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with dash at the start",
			baseName:     "-hello-World",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with underscore",
			baseName:     "hello_World",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with space",
			baseName:     "hello World",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with tab",
			baseName:     "hello\tWorld",
			expectedName: "hello-world",
		},
		{
			scenario:     "name with dot",
			baseName:     "hello.World",
			expectedName: "hello-world",
		},
		{
			scenario:     "name consecutive upper case letter",
			baseName:     "helloWORLD",
			expectedName: "hello-world",
		},
		{
			scenario:     "name consecutive upper case letter and dash",
			baseName:     "helloW-ORLD",
			expectedName: "hello-w-orld",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			name := normalizeOptionName(test.baseName, "-")
			require.Equal(t, test.expectedName, name)
		})
	}
}

func TestNormalizeCLIOptionName(t *testing.T) {
	name := normalizeCLIOptionName("hello world- i great_lolYeah")
	require.Equal(t, "hello-world-i-great-lol-yeah", name)
}

func TestNormalizeEnvOptionName(t *testing.T) {
	name := normalizeEnvOptionName("hello world- i great_lolYeah")
	require.Equal(t, "HELLO_WORLD_I_GREAT_LOL_YEAH", name)
}

func BenchmarkNormalizeOptionName(b *testing.B) {
	s := "hello world- i great_lolYeahPORIGON\tbwa"

	for b.Loop() {
		normalizeOptionName(s, "-")
	}
}
