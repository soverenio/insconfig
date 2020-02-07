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

package main

import (
	"flag"
	"strings"

	"github.com/insolar/insconfig"
)

var testflag1 = flag.String("testflag1", "", "testflag1")
var testflag2 = flag.String("testflag2", "", "testflag2")

type Config struct {
	Protocol           string
	Address            string
	FixedPublicAddress string
	HostNetwork        HostNetwork
}

type HostNetwork struct {
	MinTimeout          int
	MaxTimeout          int
	TimeoutMult         int
	SignMessages        bool
	HandshakeSessionTTL int32
}

func (c Config) GetConfig() interface{} {
	return &c
}

func main() {
	params := insconfig.Params{
		ConfigStruct: Config{},
		EnvPrefix:    "observer",
	}
	insConfigurator := insconfig.NewInsConfigurator(params, insconfig.DefaultConfigPathGetter{
		GoFlags: flag.CommandLine,
	})
	parsedConf, err := insConfigurator.Load()
	if err != nil {
		panic(err)
	}
	cfg := parsedConf.(*Config)
	insConfigurator.PrintConfig(cfg)

	if testflag1 == nil || len(*testflag1) == 0 || len(strings.TrimSpace(*testflag1)) == 0 {
		panic("testflag1 should be provided")
	}
	if testflag2 == nil || len(*testflag2) == 0 || len(strings.TrimSpace(*testflag2)) == 0 {
		panic("testflag2 should be provided")
	}
}
