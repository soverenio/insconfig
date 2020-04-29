//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package insconfig_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/insolar/insconfig"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type Level3 struct {
	Level3text string
	NullString *string
}
type Level2 struct {
	Level2text string
	Level3     Level3
}
type CfgStruct struct {
	Level1text string
	Level2     Level2
	MapField   map[string]Level2
	Map2       map[string]Level3
}

type anonymousEmbeddedStruct struct {
	CfgStruct `mapstructure:",squash"`
	Level4    string
}

type testPathGetter struct {
	Path string
}

func (g testPathGetter) GetConfigPath() string {
	return g.Path
}

func Test_Load(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, cfg.Level1text, "text1")
			require.Equal(t, cfg.Level2.Level2text, "text2")
			require.Equal(t, cfg.Level2.Level3.Level3text, "text3")
			require.Len(t, cfg.MapField, 2)
			key1 := cfg.MapField["key1"]
			require.Equal(t, key1.Level2text, "key1text2")
			require.Equal(t, key1.Level3.Level3text, "key1text3")
			require.Nil(t, key1.Level3.NullString)
			key2 := cfg.MapField["key2"]
			require.Equal(t, key2.Level2text, "key2text2")
			require.Equal(t, key2.Level3.Level3text, "key2text3")
			require.NotNil(t, key2.Level3.NullString)
		})

		t.Run("null string test", func(t *testing.T) {
			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config2.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Nil(t, cfg.Level2.Level3.NullString)
		})

		t.Run("embedded struct flatten test", func(t *testing.T) {
			cfg := anonymousEmbeddedStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config3.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, cfg.Level4, "text4")
		})

		t.Run("fail extra in yaml", func(t *testing.T) {
			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_wrong.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "nonexistent")
		})

		t.Run("fail not enough keys", func(t *testing.T) {
			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_wrong2.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)

			require.Contains(t, err.Error(), "level1text")
			require.Contains(t, err.Error(), "map2.<-key->.level3text")
			require.Contains(t, err.Error(), "map2.<-key->.nullstring")
		})

		t.Run("fail required file not found", func(t *testing.T) {
			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"nonexistent.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "nonexistent.yaml")
		})
	})

	t.Run("env", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_LEVEL1TEXT", "newTextValue1")
			_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL2TEXT", "newTextValue2")
			_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL3_LEVEL3TEXT", "newTextValue3")
			_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL3_NULLSTRING", "text")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL2TEXT", "1")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL3_LEVEL3TEXT", "2")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL3_NULLSTRING", "3")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL2TEXT", "21")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_LEVEL3TEXT", "22")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_NULLSTRING", "23")
			_ = os.Setenv("TESTPREFIX_MAP2_KEY3_LEVEL3TEXT", "32")
			_ = os.Setenv("TESTPREFIX_MAP2_KEY3_NULLSTRING", "33")
			defer os.Unsetenv("TESTPREFIX_LEVEL1TEXT")
			defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL2TEXT")
			defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL3_LEVEL3TEXT")
			defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL3_NULLSTRING")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL2TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL3_LEVEL3TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL3_NULLSTRING")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL2TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_LEVEL3TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_NULLSTRING")
			defer os.Unsetenv("TESTPREFIX_MAP2_KEY3_LEVEL3TEXT")
			defer os.Unsetenv("TESTPREFIX_MAP2_KEY3_NULLSTRING")

			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{""},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, "newTextValue1", cfg.Level1text)
			require.Equal(t, "newTextValue2", cfg.Level2.Level2text)
			require.Equal(t, "newTextValue3", cfg.Level2.Level3.Level3text)
			mapField := cfg.MapField
			require.Len(t, mapField, 2)
			require.Equal(t, "1", mapField["key1"].Level2text)
			require.Equal(t, "2", mapField["key1"].Level3.Level3text)
			require.Equal(t, "3", *mapField["key1"].Level3.NullString)
			require.Equal(t, "21", mapField["key2"].Level2text)
			require.Equal(t, "22", mapField["key2"].Level3.Level3text)
			require.Equal(t, "23", *mapField["key2"].Level3.NullString)
			map2 := cfg.Map2
			require.Len(t, map2, 1)
			require.Equal(t, "32", map2["key3"].Level3text)
			require.Equal(t, "33", *map2["key3"].NullString)
		})

		t.Run("fail not enough keys", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_LEVEL1TEXT", "newTextValue1")
			_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL3_LEVEL3TEXT", "newTextValue3")
			_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL3_NULLSTRING", "text")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL2TEXT", "1")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL3_NULLSTRING", "3")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY3_LEVEL2TEXT", "31")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL2TEXT", "21")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_LEVEL3TEXT", "22")
			_ = os.Setenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_NULLSTRING", "23")
			defer os.Unsetenv("TESTPREFIX_LEVEL1TEXT")
			defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL3_LEVEL3TEXT")
			defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL3_NULLSTRING")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL2TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY1_LEVEL3_NULLSTRING")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY3_LEVEL2TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL2TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_LEVEL3TEXT")
			defer os.Unsetenv("TESTPREFIX_MAPFIELD_KEY2_LEVEL3_NULLSTRING")

			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{""},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "level2.level2text")
			require.Contains(t, err.Error(), "mapfield.key1.level3.level3text")
			require.Contains(t, err.Error(), "mapfield.key3.level3.level3text")
			require.Contains(t, err.Error(), "mapfield.key3.level3.nullstring")
			require.Contains(t, err.Error(), "map2.<-key->.nullstring")
			require.Contains(t, err.Error(), "map2.<-key->.level3text")
		})
	})

	t.Run("yaml and env", func(t *testing.T) {
		t.Run("env overrides yaml partially", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL2TEXT", "newTextValue")
			defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL2TEXT")
			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, cfg.Level1text, "text1")
			require.Equal(t, cfg.Level2.Level2text, "newTextValue")
			require.Equal(t, cfg.Level2.Level3.Level3text, "text3")
		})

		t.Run("env adds values that is not in the yaml", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_LEVEL1TEXT", "newTextValue1")
			_ = os.Setenv("TESTPREFIX_MAP2_ONE_LEVEL3TEXT", "newTextValue1")
			_ = os.Setenv("TESTPREFIX_MAP2_ONE_NULLSTRING", "newTextValue1")
			defer os.Unsetenv("TESTPREFIX_LEVEL1TEXT")
			defer os.Unsetenv("TESTPREFIX_MAP2_ONE_LEVEL3TEXT")
			defer os.Unsetenv("TESTPREFIX_MAP2_ONE_NULLSTRING")
			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_wrong2.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, cfg.Level1text, "newTextValue1")
			require.Equal(t, cfg.Level2.Level2text, "text2")
			require.Equal(t, cfg.Level2.Level3.Level3text, "text3")
		})

		t.Run("embedded struct override by env", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL2TEXT", "newTextValue")
			defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL2TEXT")

			cfg := anonymousEmbeddedStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config3.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, cfg.Level2.Level2text, "newTextValue")
		})

		t.Run("fail extra in env", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_NONEXISTENT_VALUE", "123")
			defer os.Unsetenv("TESTPREFIX_NONEXISTENT_VALUE")

			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "nonexistent")
		})

		t.Run("fail extra in env with empty value", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_NONEXISTENT_VALUE1", "")
			_ = os.Setenv("TESTPREFIX_NONEXISTENT_VALUE2", "")
			defer os.Unsetenv("TESTPREFIX_NONEXISTENT_VALUE1")
			defer os.Unsetenv("TESTPREFIX_NONEXISTENT_VALUE2")

			cfg := CfgStruct{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config.yaml"},
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "nonexistent.value1")
			require.Contains(t, err.Error(), "nonexistent.value2")
		})
	})

	t.Run("map in config", func(t *testing.T) {
		type MapValue struct {
			Str  string
			Num  int
			Flag bool
		}

		type OneMap struct {
			One map[string]MapValue
		}

		type TwoMaps struct {
			One map[string]MapValue
			Two map[string]MapValue
		}

		t.Run("map upper level yaml", func(t *testing.T) {
			cfg := map[string]MapValue{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_map_upper_level.yaml"},
				FileNotRequired:  false,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Len(t, cfg, 2)
			first := cfg["first"]
			require.Equal(t, "first-str", first.Str)
			require.Equal(t, 1, first.Num)
			require.True(t, first.Flag)
			second := cfg["second"]
			require.Equal(t, "second-str", second.Str)
			require.Equal(t, 2, second.Num)
			require.False(t, second.Flag)
		})

		t.Run("map upper level env", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_FIRST_STR", "first-str")
			_ = os.Setenv("TESTPREFIX_FIRST_NUM", "1")
			_ = os.Setenv("TESTPREFIX_FIRST_FLAG", "true")
			_ = os.Setenv("TESTPREFIX_SECOND_STR", "second-str")
			_ = os.Setenv("TESTPREFIX_SECOND_NUM", "2")
			_ = os.Setenv("TESTPREFIX_SECOND_FLAG", "false")
			defer os.Unsetenv("TESTPREFIX_FIRST_STR")
			defer os.Unsetenv("TESTPREFIX_FIRST_NUM")
			defer os.Unsetenv("TESTPREFIX_FIRST_FLAG")
			defer os.Unsetenv("TESTPREFIX_SECOND_STR")
			defer os.Unsetenv("TESTPREFIX_SECOND_NUM")
			defer os.Unsetenv("TESTPREFIX_SECOND_FLAG")

			cfg := map[string]MapValue{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{""},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Len(t, cfg, 2)
			first := cfg["first"]
			require.Equal(t, "first-str", first.Str)
			require.Equal(t, 1, first.Num)
			require.True(t, first.Flag)
			second := cfg["second"]
			require.Equal(t, "second-str", second.Str)
			require.Equal(t, 2, second.Num)
			require.False(t, second.Flag)
		})

		t.Run("one map", func(t *testing.T) {
			cfg := OneMap{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_one_map.yaml"},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Len(t, cfg.One, 2)
			first := cfg.One["first"]
			require.Equal(t, "first-str", first.Str)
			require.Equal(t, 1, first.Num)
			require.True(t, first.Flag)
			second := cfg.One["second"]
			require.Equal(t, "second-str", second.Str)
			require.Equal(t, 2, second.Num)
			require.False(t, second.Flag)
		})

		t.Run("two maps one level", func(t *testing.T) {
			cfg := TwoMaps{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_two_maps.yaml"},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Len(t, cfg.One, 2)
			first := cfg.One["first"]
			require.Equal(t, "first-str", first.Str)
			require.Equal(t, 1, first.Num)
			require.True(t, first.Flag)
			second := cfg.One["second"]
			require.Equal(t, "second-str", second.Str)
			require.Equal(t, 2, second.Num)
			require.False(t, second.Flag)
			require.Len(t, cfg.Two, 1)
			firstOfTwo := cfg.Two["first"]
			require.Equal(t, "two-first-str", firstOfTwo.Str)
			require.Equal(t, 3, firstOfTwo.Num)
			require.True(t, firstOfTwo.Flag)
		})

		t.Run("map string values yaml", func(t *testing.T) {
			cfg := make(map[string]string, 0)
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_map_str.yaml"},
				FileNotRequired:  false,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, "first-str", cfg["str"])
			require.Equal(t, "1", cfg["num"])
		})

		t.Run("map string values env", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_STR", "first-str")
			_ = os.Setenv("TESTPREFIX_NUM", "1")
			defer os.Unsetenv("TESTPREFIX_STR")
			defer os.Unsetenv("TESTPREFIX_NUM")

			cfg := make(map[string]string, 0)
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{""},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.NoError(t, err)
			require.Equal(t, "first-str", cfg["str"])
			require.Equal(t, "1", cfg["num"])
		})

		t.Run("fail two nested maps yaml", func(t *testing.T) {
			type TwoNested struct {
				One map[string]map[string]MapValue
			}

			cfg := TwoNested{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_two_maps_nested.yaml"},
				FileNotRequired:  false,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Equal(t, "nested maps are not allowed in config", err.Error())
		})

		t.Run("fail two nested maps env", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_ONE_FIRST_FIRST_STR", "first-str")
			_ = os.Setenv("TESTPREFIX_ONE_FIRST_FIRST_NUM", "1")
			_ = os.Setenv("TESTPREFIX_ONE_FIRST_FIRST_FLAG", "true")
			defer os.Unsetenv("TESTPREFIX_ONE_FIRST_FIRST_STR")
			defer os.Unsetenv("TESTPREFIX_ONE_FIRST_FIRST_NUM")
			defer os.Unsetenv("TESTPREFIX_ONE_FIRST_FIRST_FLAG")

			type TwoNested struct {
				One map[string]map[string]MapValue
			}

			cfg := TwoNested{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{""},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Equal(t, "nested maps are not allowed in config", err.Error())
		})

		t.Run("fail map-struct-map yaml", func(t *testing.T) {
			type StructMap struct {
				Two map[string]MapValue
			}
			type MapStructMap struct {
				One map[string]StructMap
			}

			cfg := MapStructMap{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_map_struct_map.yaml"},
				FileNotRequired:  false,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Equal(t, "nested maps are not allowed in config", err.Error())
		})

		t.Run("fail map-struct-map env", func(t *testing.T) {
			_ = os.Setenv("TESTPREFIX_ONE_FIRST_TWO_FIRST_STR", "first-first")
			_ = os.Setenv("TESTPREFIX_ONE_FIRST_TWO_FIRST_NUM", "1")
			_ = os.Setenv("TESTPREFIX_ONE_FIRST_TWO_FIRST_FLAG", "true")
			defer os.Unsetenv("TESTPREFIX_ONE_FIRST_TWO_FIRST_STR")
			defer os.Unsetenv("TESTPREFIX_ONE_FIRST_TWO_FIRST_NUM")
			defer os.Unsetenv("TESTPREFIX_ONE_FIRST_TWO_FIRST_FLAG")

			type StructMap struct {
				Two map[string]MapValue
			}
			type MapStructMap struct {
				One map[string]StructMap
			}

			cfg := MapStructMap{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{""},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Equal(t, "nested maps are not allowed in config", err.Error())
		})

		t.Run("fail one map int key", func(t *testing.T) {
			type MapIntKey struct {
				One map[int]MapValue
			}
			cfg := MapIntKey{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_map_int_key.yaml"},
				FileNotRequired:  false,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "maps in config must have string keys but got:")
		})

		t.Run("fail one map struct key", func(t *testing.T) {
			type MapIntKey struct {
				One map[MapValue]MapValue
			}
			cfg := MapIntKey{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{""},
				FileNotRequired:  true,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "maps in config must have string keys but got:")
		})

		t.Run("fail key duplication", func(t *testing.T) {
			cfg := OneMap{}
			params := insconfig.Params{
				EnvPrefix:        "testprefix",
				ConfigPathGetter: testPathGetter{"testdata/test_config_key_duplication.yaml"},
				FileNotRequired:  false,
			}

			insConfigurator := insconfig.New(params)
			err := insConfigurator.Load(&cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "failed to unmarshal config file into configuration structure")
			require.Contains(t, err.Error(), `key "first" already set in map`)
		})
	})
}

type Y struct {
	F int `insconfig:"111| the F comment"`
}

type X struct {
	A string `insconfig:"Adefault|large comment A with pipe='|'"`
	B string `insconfig:"Bdefault|large comment B with pipe='|'"`
	E *Y     `insconfig:"|---------------------------" yaml:"sacsacasc"`
	C int
	D uint8
	G map[string]int
	H []int
}

func Test_TemplateTo(t *testing.T) {
	x := &X{"1", "2", &Y{}, 3, 4, map[string]int{}, []int{}}
	w := &bytes.Buffer{}
	err := insconfig.NewYamlTemplater(x).TemplateTo(w)
	require.NoError(t, err)
	s := w.String()

	nx := X{}
	require.NoError(t, yaml.Unmarshal(w.Bytes(), &nx))
	require.NotNil(t, nx.E)

	require.Contains(t, s, "#large comment A with pipe='|'")
	require.Contains(t, s, "a: Adefault # string")
	require.Contains(t, s, "#large comment B with pipe='|'")
	require.Contains(t, s, "b: Bdefault # string")
	require.Contains(t, s, "#---------------------------")
	require.Contains(t, s, "sacsacasc:")
	require.Contains(t, s, "# the F comment")
	require.Contains(t, s, "f: 111 # int")
	require.Contains(t, s, "c:  # int")
	require.Contains(t, s, "d:  # uint8")
	require.Contains(t, s, "g: # <map> of int")
	require.Contains(t, s, "h: # <array> of int")

}

type Z struct {
	A A
}

type A struct{}

func (A) TemplateTo(w io.Writer, m *insconfig.YamlTemplater) error {
	return errors.New("")
}

func Test_FailTemplateTo(t *testing.T) {
	w := &bytes.Buffer{}
	require.NotNil(t, insconfig.NewYamlTemplater(Z{}).TemplateTo(w))
}
