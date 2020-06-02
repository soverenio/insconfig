# Insconfig

![test](https://github.com/insolar/insconfig/workflows/test/badge.svg)

Configs are used for unified approach for all Insolar applications (parameters, configs; assured ledger, observer, mainnet, generic Insolar explorer,  )

Insolar is extensively documented. If you need more information on what is Insolar, visit Insolar Docs.

Config management library.
This is the wrapper on https://github.com/spf13/viper library

## Key features
- .yaml format of config files
- No default config paths. A path should be explicitly set by the "--config/-c" flag. Optionally you can override this by implementing ConfigPathGetter (look at `configurator_test.go`)
- The environment values override file values
- Possible to use only ENV, without a configuration file at all
- Optional printing of the config file content to the log file at an app launch. 
- Hiding sensitive data via tags
- No default values for a configuration file. All values should be set explicitly, otherwise insonfig returns an error
- No unnecessary values (both in a configuration file and ENV), otherwise insonfig returns an error
- Supports custom flags, go flags and pflags
- Doesn't support overriding configutation files by flags
- Generates an empty .yaml file with field descriptions
- By default, insconfig adds the "--config" flag
- [work in progress] By default, insconfig adds the "--gen-config" flag
- Doesn't support overriding app configuration on the fly
- Supports custom viper decode hooks

## Running example 

```
go run ./example/example.go --config="./example/example_config.yaml"
```

## Usage

If you don't use any flags and DefaultPathGetter (it adds the --config flag and you can start away).

```
mconf := &Config{}	
	params := insconfig.Params{
		EnvPrefix:        "example",
		ConfigPathGetter: &insconfig.DefaultPathGetter{},
	}
	insConfigurator := insconfig.New(params)
	if err := insConfigurator.Load(cfg); err != nil {
		panic(err)
	}
	insConfigurator.ToYaml(cfg)
```


If you want to manage flags yourself.

Don't forget to add `github.com/insolar/insconfig` to the import section.

With custom go flags (from example.go)
```
go
    var flag_example_1 = flag.String("flag_example_1", "", "flag_example_1_desc")
    mconf := Config{}
    params := insconfig.Params{
	EnvPrefix:    "example",
	ConfigPathGetter: &insconfig.FlagPathGetter{
		GoFlags: flag.CommandLine,
	},
    }
    insConfigurator := insconfig.New(params)
    _ = insConfigurator.Load(&mconf)
    fmt.Println(flag_example_1)
```

With custom [spf13/pflags](https://github.com/spf13/pflag)
```go
    var flag_example_1 = pflag.String("flag_example_1", "", "flag_example_1_desc")
    mconf := Config{}
    params := insconfig.Params{
        EnvPrefix:    "example",
        ConfigPathGetter: &insconfig.PFlagPathGetter{
            PFlags: pflag.CommandLine,
        },
    }
    insConfigurator := insconfig.New(params)
    _ = insConfigurator.Load(&mconf)
    fmt.Println(testflag1)
```

With [spf13/cobra](https://github.com/spf13/cobra). Cobra doesn't provide tools to manage flags parsing, so you need to add the "--config" flag yourself.

```go
func main () {
    var configPath string
    rootCmd := &cobra.Command{
        Use: "insolard",
    }
    rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file")
    _ = rootCmd.MarkPersistentFlagRequired("config")
    err := rootCmd.Execute()

    // ...

    // To set your path from flag to insconfig you need to implement simple ConfigPathGetter interface and return path 
    type stringPathGetter struct {
        Path string
    }
    
    func (g *stringPathGetter) GetConfigPath() string {
        return g.Path
    }
}

func read(){
    mconf := Config{}
    params := insconfig.Params{
        EnvPrefix:        "example",
        ConfigPathGetter: &stringPathGetter{Path: configPath},
        FileRequired:     false,
    }
    insConfigurator := insconfig.NewInsConfigurator(h.Params)
    err := insConfigurator.Load(&mconf)
    println(insconfig.ToString(mconf))
}
```


### Create a configuration template
If you want to get a config file example and you have a ready config structure, you can use this code as a reference example:

```go
    type Config struct {
    ... 
        Field Type `insonfig:"default_value|Commentary for this field"`
    }
    ...
    insconfig.NewYamlTemplater(new(Config)).TemplateTo(os.StdOut)
```

Tip: You can use tags to enrich a field with a default value and a comment for this value; both will end up in your template. 

### Create a configuration template with hidden fields

If you have some sensitive data you may want to hide it in a config. You can use the `insconfigsecret:` tag to hide this data.

```go
    type Config struct {
    ...
	    Pass string `insconfigsecret:""`
    }
    ...
    insconfig.NewYamlDumper(Config).DumpTo(os.StdOut)
```

### Using maps in a configuration file

You can use maps in a configuration file, althought with some limitations:
- Only String type keys are allowed
- A map cannot be nested and has to be used on the first level
- Nested maps (directly or in a struct) are not allowed

