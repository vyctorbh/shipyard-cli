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
	app.Name = "Shipyard CLI"
	app.Version = "0.1.0"
	app.Usage = "Command line interface for Shipyard"
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
			Name:      "apps",
			ShortName: "",
			Usage:     "Application Management",
			Flags: []cli.Flag{
				cli.StringFlag{"action, a", "show", "Show Applications"},
				cli.StringFlag{"name, n", "", "Application Name (optional)"},
			},
			Action: func(c *cli.Context) {
				applicationsAction(c)
			},
		},
		{
			Name:      "containers",
			ShortName: "",
			Usage:     "Container Management",
			Flags: []cli.Flag{
                                cli.BoolFlag{"start", "Start Container"},
                                cli.BoolFlag{"stop", "Stop Container"},
                                cli.BoolFlag{"restart", "Restart Container"},
                                cli.BoolFlag{"remove", "Remove Container"},
				cli.StringFlag{"id, i", "", "Container ID (optional)"},
				cli.BoolFlag{"all", "Show all containers (optional)"},
			},
			Action: func(c *cli.Context) {
				containersAction(c)
			},
		},
		{
			Name:      "images",
			ShortName: "",
			Usage:     "Image Management",
			Flags: []cli.Flag{
				cli.StringFlag{"id, i", "", "ID of Image (optional)"},
			},
			Action: func(c *cli.Context) {
				imagesAction(c)
			},
		},
		{
			Name:      "hosts",
			ShortName: "",
			Usage:     "Host Management",
			Flags: []cli.Flag{
				cli.StringFlag{"name, n", "", "Name of Host (optional)"},
			},
			Action: func(c *cli.Context) {
				hostsAction(c)
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
