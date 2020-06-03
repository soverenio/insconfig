[<img src="https://github.com/insolar/doc-pics/raw/master/st/github-readme-banner.png">](http://insolar.io/?utm_source=Github)

# Insolar configurations

![test](https://github.com/insolar/insconfig/workflows/test/badge.svg)

Insolar configurations unify configuring for all Insolar applications: [Assured Ledger](https://github.com/insolar/assured-ledger), [Insolar Observer](https://github.com/insolar/observer), [Insolar MainNet](https://github.com/insolar/mainnet), [Insolar Explorer](https://github.com/insolar/block-explorer) and more.

Configurations are consolidated into a **single configuration library** that is a wrapper for [Viper](https://github.com/spf13/viper).

Configurations are used by Insolar applications as a dependency.

Insolar is extensively documented. If you need more information on what is Insolar, visit [Insolar Docs](http://docs.insolar.io/quick_overview.html).

## Key features
- .yaml format of configuration files.
- No default configuration paths. A path should be explicitly set by the "--config/-c" flag. Optionally, you can override this by implementing ConfigPathGetter (look at `configurator_test.go` for details).
- Environment values override file values.
- Option to use only ENV, without a configuration file at all.
- Option to write the config file content to the log file at an app launch. 
- Hiding sensitive data via the `insconfigsecret` tag.
- No default values in a configuration file. All values should be set explicitly, otherwise the library returns an error.
- No unnecessary field or parameters both in a configuration file and ENV, otherwise the library returns an error. Consider as unecessary: fields in a config struct unused in a configuration file, old or obsolete parameters in a configuration file that are not currently used, unused parameters in ENV.
- Support of custom flags, go flags and pflags.
- No overriding configutation files by flags.
- Option to generate an empty .yaml file with field descriptions.
- Automatic adding of the `--config` flag
- [work in progress] Automatic of the `--gen-config` flag
- No overriding app configuration on the fly.
- Support of custom Viper decode hooks.

## Usage

### In terminal

Consider this example:

```
go run ./example/example.go --config="./example/example_config.yaml"
```

## In your code

Tip: Don't forget to add `github.com/insolar/insconfig` to the import section.

### No flags

If you don't use any flags and DefaultPathGetter, which adds the `--config` flag so you can start right away, consider this example:

```
mconf := &Config{}	
	params := insconfig.Params{
		EnvPrefix:        "example",
		ConfigPathGetter: &insconfig.DefaultPathGetter{},
	}
	insConfigurator := insconfig.New(params)
	if err := insConfigurator.Load(mconf); err != nil {
		panic(err)
	}
	insConfigurator.ToYaml(mconf)
```

### Custom flags

#### Custom Go flags (from example.go)

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

#### Custom [spf13/pflags](https://github.com/spf13/pflag)

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

#### Custom [spf13/cobra flags](https://github.com/spf13/cobra). 

Note: Cobra doesn't provide tools for managing flags parsing, so you need to add the `--config flag yourself.

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


### Generating a configuration template

Tip: You can use tags to enrich a field with a default value and a comment for this value; both will end up in your template. 

If you want to create a config file example and you have a ready config structure, consider this example:

```go
    type Config struct {
    ... 
        Field Type `insonfig:"default_value|Commentary for this field"`
    }
    ...
    insconfig.NewYamlTemplater(new(Config)).TemplateTo(os.StdOut)
```

### Generating a configuration template with hidden fields

If you have some sensitive data you may want to hide it in a config. You can use the `insconfigsecret` tag to hide such data.

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

## Contribute!

Feel free to submit issues, fork the repository and send pull requests! 

To make the process smooth for both reviewers and contributors, familiarize yourself with the following guidelines:

1. [Open source contributor guide](https://github.com/freeCodeCamp/how-to-contribute-to-open-source).
2. [Style guide: Effective Go](https://golang.org/doc/effective_go.html).
3. [List of shorthands for Go code review comments](https://github.com/golang/go/wiki/CodeReviewComments).

When submitting an issue, **include a complete test function** that reproduces it.

Thank you for your intention to contribute to the Insolar Mainnet project. As a company developing open-source code, we highly appreciate external contributions to our project.

## Contacts

If you have any additional questions, join our [developers chat on Telegram](https://t.me/InsolarTech).

Our social media:

[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-facebook.png" width="36" height="36">](https://facebook.com/insolario)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-twitter.png" width="36" height="36">](https://twitter.com/insolario)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-medium.png" width="36" height="36">](https://medium.com/insolar)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-youtube.png" width="36" height="36">](https://youtube.com/insolar)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-reddit.png" width="36" height="36">](https://www.reddit.com/r/insolar/)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-linkedin.png" width="36" height="36">](https://www.linkedin.com/company/insolario/)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-instagram.png" width="36" height="36">](https://instagram.com/insolario)
[<img src="https://github.com/insolar/doc-pics/raw/master/st/ico-social-telegram.png" width="36" height="36">](https://t.me/InsolarAnnouncements) 

## License

This project is licensed under the terms of the [Insolar License 1.0](LICENSE.md).
