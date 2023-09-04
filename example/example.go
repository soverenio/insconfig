package main

import (
	"flag"
	"fmt"

	"github.com/soverenio/insconfig"
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
