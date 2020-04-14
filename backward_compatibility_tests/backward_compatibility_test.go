// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
