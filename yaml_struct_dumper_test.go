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
	"fmt"
	"testing"

	"github.com/insolar/insconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Simple       string            `insconfig:"Example secret value" insconfigsecret:""`
	List         []string          `insconfig:"List"`
	Map          map[string]string `insconfig:"Map"`
	ListOfStruct []Inner           `insconfig:"List of structs"`
	MapOfStruct  map[string]Inner  `insconfig:"Map of structs"`
	Inner        Inner
}

type Inner struct {
	F1 string `insconfig:"First inner field"`
	F2 string `insconfig:"Second inner field"`
}

func NewConfig() Config {
	cfg := Config{
		Simple: "example",
		List:   []string{"val1", "val2", "val3"},
		Map: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
		ListOfStruct: []Inner{
			{
				F1: "firstValue",
				F2: "secondValue",
			}, {
				F1: "value1",
				F2: "value2",
			},
		},
		MapOfStruct: map[string]Inner{
			"key1": {
				F1: "value1",
				F2: "value2",
			},
			"key2": {
				F1: "value3",
				F2: "value4",
			},
		},
		Inner: Inner{
			F1: "innerField1",
			F2: "innerField2",
		},
	}

	return cfg
}

func TestTemplaterStruct_TemplateTo(t *testing.T) {
	w := &bytes.Buffer{}
	err := insconfig.NewYamlTemplaterStruct(NewConfig()).TemplateTo(w)
	require.NoError(t, err)
	s := w.String()
	fmt.Println(s)

	assert.Contains(t, s, `# Example secret value
simple: example # string
# List
list: # <array> of string`)
	assert.Contains(t, s, `- val1 # string`)
	assert.Contains(t, s, `- val2 # string`)
	assert.Contains(t, s, `- val3 # string`)
	assert.Contains(t, s, `
# Map
map: # <map> of string`)
	assert.Contains(t, s, `
  key1: value1 # string`)
	assert.Contains(t, s, `
  key2: value2 # string`)
	assert.Contains(t, s, `
  key3: value3 # string`)
	assert.Contains(t, s, `
# List of structs
listofstruct: # <array> of insconfig_test.Inner`)
	assert.Contains(t, s, `
  -
    # First inner field
    f1: firstValue # string
    # Second inner field
    f2: secondValue # string`)
	assert.Contains(t, s, `
  -
    # First inner field
    f1: value1 # string
    # Second inner field
    f2: value2 # string`)
	assert.Contains(t, s, `
# Map of structs
mapofstruct: # <map> of insconfig_test.Inner`)
	assert.Contains(t, s, `
  key1:
    # First inner field
    f1: value1 # string
    # Second inner field
    f2: value2 # string`)
	assert.Contains(t, s, `
  key2:
    # First inner field
    f1: value3 # string
    # Second inner field
    f2: value4 # string`)
	assert.Contains(t, s, `
inner:
  # First inner field
  f1: innerField1 # string
  # Second inner field
  f2: innerField2 # string`)
}
