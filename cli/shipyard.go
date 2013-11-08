package main

import (
    "github.com/ehazlett/shipyard-go/shipyard"
    "github.com/codegangsta/cli"
    "fmt"
    "flag"
    "os"
    "path"
    "encoding/json"
    "io/ioutil"
)

var (
    apiUsername string
    apiKey string
    apiUrl string
    apiVersion string
)

type Config struct {
    Username    string
    ApiKey      string
    Url         string
    Version     string
}

func loadConfig() (Config, error) {
    configFile := path.Join(os.Getenv("HOME"), ".shipyard.cfg")
    var config = Config{}
    _, err := os.Stat(configFile)
    if err != nil || apiUrl != "" || apiUsername != "" || apiKey != "" {
        if apiUrl == "" {
            fmt.Println("Error: You must specify a Shipyard URL")
            flag.PrintDefaults()
            os.Exit(1)
        }
        if apiUsername == "" {
            fmt.Println("Error: You must specify a Shipyard Username")
            os.Exit(1)
        }
        if apiKey == "" {
            fmt.Println("Error: You must specify a Shipyard API Key")
            os.Exit(1)
        }
        // config
        config = Config{
            Username: apiUsername,
            ApiKey: apiKey,
            Url: apiUrl,
            Version: apiVersion,
        }
        // save config
        cfg, err := json.Marshal(config)
        b := []byte(cfg)
        err = ioutil.WriteFile(configFile, b, 0600)
        if err != nil { panic(err) }
    } else {
        cfg, err := ioutil.ReadFile(configFile)
        if err != nil { panic(err) }
        b := []byte(cfg)
        config = Config{}
        err = json.Unmarshal(b, &config)
        if err != nil { panic(err) }
    }
    return config, nil
}

func getAPI() shipyard.API {
    config, _ := loadConfig()
    api := shipyard.NewAPI(config.Username, config.ApiKey, config.Url, config.Version)
    return *api
}

func main() {
    app := cli.NewApp()
    app.Name = "Shipyard"
    app.Version = "0.1.0"
    app.Usage = "Shipyard CLI"
    app.Flags = []cli.Flag {
        cli.StringFlag{"username", "", "Shipyard API Username"},
        cli.StringFlag{"key", "", "Shipyard API Key"},
        cli.StringFlag{"url", "", "Shipyard URL"},
        cli.StringFlag{"api-version", "1", "Shipyard API Version"},
    }
    app.Action = func(c *cli.Context) {
        apiUsername = c.String("username")
        apiKey = c.String("key")
        apiUrl = c.String("url")
        apiVersion = c.String("api-version")
        if len(c.Args()) == 0 {
            cli.ShowAppHelp(c)
            os.Exit(2)
        }
        // check for valid command
        validCommand := false
        command := c.Args()[0]
        for _,v := range app.Commands {
            if v.Name == command || v.ShortName == command { validCommand = true }
        }
        if ! validCommand {
            cli.ShowAppHelp(c)
        }
    }
    app.Commands = []cli.Command{
        {
            Name:      "applications",
            ShortName: "apps",
            Usage:     "Manage Applications",
            Action: func(c *cli.Context) {
            },
        },
        {
            Name:      "containers",
            ShortName: "cnt",
            Usage:     "Manage Containers",
            Action: func(c *cli.Context) {
            },
        },
        {
            Name:      "hosts",
            ShortName: "hosts",
            Usage:     "Manage Hosts",
            Action: func(c *cli.Context) {
            },
        },
        {
            Name:      "config",
            ShortName: "cfg",
            Usage:     "Show current Shipyard config",
            Action: func(c *cli.Context) {
                config, _ := loadConfig()
                LogMessage(fmt.Sprintf("Endpoint: %v", config.Url))
                LogMessage(fmt.Sprintf("Username: %v", config.Username))
                LogMessage(fmt.Sprintf("Version: %v", config.Version))
                LogMessage(fmt.Sprintf("APIKey: %v...", config.ApiKey[0:5]))
            },
        },
        {
            Name:      "info",
            ShortName: "info",
            Usage:     "Show Shipyard Info",
            Action: func(c *cli.Context) {
                api := getAPI()
                fmt.Println(api.GetInfo())
            },
        },
    }
    // run cli
    app.Run(os.Args)
}
