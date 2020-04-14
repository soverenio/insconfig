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
	"fmt"

	"github.com/insolar/insconfig"
)

var testflag1 = flag.String("testflag1", "", "testflag1")
var testflag2 = flag.String("testflag2", "", "testflag2")

type Config struct {
	Protocol           string
	Address            string
	FixedPublicAddress string
	HostNetwork        HostNetwork
	Clients            map[string]Client
}

type HostNetwork struct {
	MinTimeout          int
	MaxTimeout          int
	TimeoutMult         int
	SignMessages        bool
	HandshakeSessionTTL int32
}

type Client struct {
	ID      int
	Address string
}

func main() {
	mconf := Config{}
	params := insconfig.Params{
		EnvPrefix: "example",
		ConfigPathGetter: &insconfig.FlagPathGetter{
			GoFlags: flag.CommandLine,
		},
	}
	insConfigurator := insconfig.New(params)
	if err := insConfigurator.Load(&mconf); err != nil {
		panic(err)
	}
	fmt.Println(*testflag1)
	fmt.Println(insConfigurator.ToYaml(mconf))
}
