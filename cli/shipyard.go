package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/shipyard/shipyard-go/shipyard"
	"io/ioutil"
	"os"
	"path"
)

var configFile = path.Join(os.Getenv("HOME"), ".shipyard.cfg")

var (
	apiUsername string
	apiKey      string
	apiUrl      string
	apiVersion  string
)

type Config struct {
	Username string
	ApiKey   string
	Url      string
	Version  string
}

func saveConfig(username string, apiKey string, url string, version string) {
	// config
	var config = Config{}
	config = Config{
		Username: username,
		ApiKey:   apiKey,
		Url:      url,
		Version:  version,
	}
	// save config
	cfg, err := json.Marshal(config)
	b := []byte(cfg)
	err = ioutil.WriteFile(configFile, b, 0600)
	if err != nil {
		panic(err)
	}
}

func loadConfig(c *cli.Context) (Config, error) {
	apiUsername = c.GlobalString("username")
	apiKey = c.GlobalString("key")
	apiUrl = c.GlobalString("url")
	apiVersion = c.GlobalString("api-version")
	var config = Config{}
	_, err := os.Stat(configFile)
	if err != nil {
		saveConfig(apiUsername, apiKey, apiUrl, apiVersion)
	} else {
		cfg, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		b := []byte(cfg)
		config = Config{}
		err = json.Unmarshal(b, &config)
		if err != nil {
			panic(err)
		}
	}
	return config, nil
}

func getAPI(c *cli.Context) shipyard.API {
	config, _ := loadConfig(c)
	api := shipyard.NewAPI(config.Username, config.ApiKey, config.Url, config.Version)
	return *api
}

func main() {
	app := cli.NewApp()
	app.Name = "Shipyard"
	app.Version = "0.1.0"
	app.Usage = "Shipyard CLI"
	app.Flags = []cli.Flag{
		cli.StringFlag{"username", "", "Shipyard API Username"},
		cli.StringFlag{"key", "", "Shipyard API Key"},
		cli.StringFlag{"url", "", "Shipyard URL"},
		cli.StringFlag{"api-version", "1", "Shipyard API Version"},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(2)
		}
		// check for valid command
		validCommand := false
		command := c.Args()[0]
		for _, v := range app.Commands {
			if v.Name == command || v.ShortName == command {
				validCommand = true
			}
		}
		if !validCommand {
			cli.ShowAppHelp(c)
		}
	}
	app.Commands = []cli.Command{
		{
			Name:      "login",
			ShortName: "",
			Usage:     "Login",
			Action: func(c *cli.Context) {
				loginAction(c)
			},
		},
		{
			Name:      "show-applications",
			ShortName: "",
			Usage:     "Show Applications",
			Flags: []cli.Flag{
				cli.StringFlag{"name", "", "Application Name (optional)"},
			},
			Action: func(c *cli.Context) {
				showApplicationsAction(c)
			},
		},
		{
			Name:      "show-containers",
			ShortName: "",
			Usage:     "Show Containers",
			Flags: []cli.Flag{
				cli.StringFlag{"id", "", "Container ID (optional)"},
			},
			Action: func(c *cli.Context) {
				showContainersAction(c)
			},
		},
		{
			Name:      "show-hosts",
			ShortName: "",
			Usage:     "Show Hosts",
			Flags: []cli.Flag{
				cli.StringFlag{"name", "", "Name of Host (optional)"},
			},
			Action: func(c *cli.Context) {
				showHostsAction(c)
			},
		},
		{
			Name:      "config",
			ShortName: "cfg",
			Usage:     "Show current Shipyard config",
			Action: func(c *cli.Context) {
				configAction(c)
			},
		},
		{
			Name:      "info",
			ShortName: "info",
			Usage:     "Show Shipyard Info",
			Action: func(c *cli.Context) {
				api := getAPI(c)
				fmt.Println(api.GetInfo())
			},
		},
	}
	// run cli
	app.Run(os.Args)
}
