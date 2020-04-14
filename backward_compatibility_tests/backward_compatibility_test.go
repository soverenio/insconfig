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

package backward_compatibility_tests

// The purpose of this package is to implement tests that refer to all existing public objects
// and functions in 'insconfig'. Once the object/function with its parameters is implemented in the library,
// means that it might be used anywhere else in its current form. 'Insconfig' must guarantee that function signature
// must remain the same and/or struct parameters have the same type as before and were not removed.

import (
	"io/ioutil"
	"testing"

	"github.com/insolar/insconfig"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

type Level4 struct {
	Level4text string
}

type Level3 struct {
	Level3text string
}

type Level2 struct {
	Level2text string
	Level3     Level3
}

type CfgStruct struct {
	Level1text string
	Level2     Level2
	Level4     Level4
}

type testPathGetter struct {
	Path string
}

func (g testPathGetter) GetConfigPath() string {
	return g.Path
}

// This test checks if the main functionality below 'insconfig' tag v0.2 has not been changed.
// The test must not be edited in the upcoming releases in order to guarantee backward compatibility
// of the 'insconfig' library.
func TestBackwardCompatibility_v02(t *testing.T) {
	cfgFileName := "test_config_ok.yaml"
	cfg := CfgStruct{}
	var f mapstructure.DecodeHookFunc = mapstructure.WeaklyTypedHook
	params := insconfig.Params{
		EnvPrefix:        "testprefix",
		ViperHooks:       []mapstructure.DecodeHookFunc{f},
		ConfigPathGetter: testPathGetter{cfgFileName},
		FileNotRequired:  false,
	}
	insConfigurator := insconfig.New(params)

	t.Run("Load", func(t *testing.T) {
		err := insConfigurator.Load(&cfg)
		require.NoError(t, err)
		require.Equal(t, cfg.Level1text, "text1")
		require.Equal(t, cfg.Level2.Level2text, "text2")
		require.Equal(t, cfg.Level2.Level3.Level3text, "text3")
		require.Equal(t, cfg.Level4.Level4text, "text5")
	})

	t.Run("ToYaml", func(t *testing.T) {
		yaml := insConfigurator.ToYaml(&cfg)
		b, _ := ioutil.ReadFile(cfgFileName)
		require.Equal(t, string(b), yaml)
	})
}
