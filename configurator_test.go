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
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insconfig"
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
	t.Run("happy", func(t *testing.T) {
		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.NoError(t, err)
		require.Equal(t, cfg.Level1text, "text1")
		require.Equal(t, cfg.Level2.Level2text, "text2")
		require.Equal(t, cfg.Level2.Level3.Level3text, "text3")
	})

	t.Run("ENV overriding", func(t *testing.T) {
		_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL2TEXT", "newTextValue")
		defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL2TEXT")
		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.NoError(t, err)
		require.Equal(t, cfg.Level1text, "text1")
		require.Equal(t, cfg.Level2.Level2text, "newTextValue")
		require.Equal(t, cfg.Level2.Level3.Level3text, "text3")
	})

	t.Run("ENV has values, that is not in config, but it should", func(t *testing.T) {
		_ = os.Setenv("TESTPREFIX_LEVEL1TEXT", "newTextValue1")
		defer os.Unsetenv("TESTPREFIX_LEVEL1TEXT")
		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config_wrong2.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.NoError(t, err)
		require.Equal(t, cfg.Level1text, "newTextValue1")
		require.Equal(t, cfg.Level2.Level2text, "text2")
		require.Equal(t, cfg.Level2.Level3.Level3text, "text3")
	})

	t.Run("ENV only, no config files", func(t *testing.T) {
		_ = os.Setenv("TESTPREFIX_LEVEL1TEXT", "newTextValue1")
		_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL2TEXT", "newTextValue2")
		_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL3_LEVEL3TEXT", "newTextValue3")
		_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL3_NULLSTRING", "text")
		defer os.Unsetenv("TESTPREFIX_LEVEL1TEXT")
		defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL2TEXT")
		defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL3_LEVEL3TEXT")
		defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL3_NULLSTRING")

		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{""},
			FileNotRequired:  true,
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.NoError(t, err)
		require.Equal(t, cfg.Level1text, "newTextValue1")
		require.Equal(t, cfg.Level2.Level2text, "newTextValue2")
		require.Equal(t, cfg.Level2.Level3.Level3text, "newTextValue3")
	})

	t.Run("extra env fail", func(t *testing.T) {
		_ = os.Setenv("TESTPREFIX_NONEXISTENT_VALUE", "123")
		defer os.Unsetenv("TESTPREFIX_NONEXISTENT_VALUE")

		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nonexistent")
	})

	t.Run("extra in file fail", func(t *testing.T) {
		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config_wrong.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nonexistent")
	})

	t.Run("not set in file fail", func(t *testing.T) {
		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config_wrong2.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Level1text")
	})

	t.Run("required file not found", func(t *testing.T) {
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

	t.Run("null string test", func(t *testing.T) {
		cfg := CfgStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config2.yaml"},
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
			ConfigPathGetter: testPathGetter{"test_config3.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.NoError(t, err)
		require.Equal(t, cfg.Level4, "text4")
	})

	t.Run("embedded struct override by env", func(t *testing.T) {
		_ = os.Setenv("TESTPREFIX_LEVEL2_LEVEL2TEXT", "newTextValue")
		defer os.Unsetenv("TESTPREFIX_LEVEL2_LEVEL2TEXT")

		cfg := anonymousEmbeddedStruct{}
		params := insconfig.Params{
			EnvPrefix:        "testprefix",
			ConfigPathGetter: testPathGetter{"test_config3.yaml"},
		}

		insConfigurator := insconfig.New(params)
		err := insConfigurator.Load(&cfg)
		require.NoError(t, err)
		require.Equal(t, cfg.Level2.Level2text, "newTextValue")
	})
}
