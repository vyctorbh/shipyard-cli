package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gcmurphy/getpass"
	"strings"
)

func loginAction(c *cli.Context) {
	api := getAPI(c)
	var url string
	var version string
	var username string
	var password string
	fmt.Printf("URL: ")
	_, urlErr := fmt.Scanln(&url)
	if urlErr != nil {
		panic(urlErr)
	}
	fmt.Printf("Username: ")
	_, usernameErr := fmt.Scanln(&username)
	if usernameErr != nil {
		panic(usernameErr)
	}
	password, passwordErr := getpass.GetPass()
	if passwordErr != nil {
		panic(passwordErr)
	}
	fmt.Printf("Version (default: 1): ")
	_, verErr := fmt.Scanln(&version)
	if verErr != nil {
		version = "1"
	}
	userData, loginErr := api.Login(url, username, password)
	if loginErr != nil {
		panic(loginErr)
	}
	saveConfig(username, userData.ApiKey, url, version)
	LogMessage("Login successful", "g")
}

func showApplicationsAction(c *cli.Context) {
	api := getAPI(c)
	appMeta, err := api.GetApplications()
	if err != nil {
		panic(err)
	}
	apps := appMeta.Objects
	// check for app op
	appName := c.String("name")
	if appName != "" {
		// check for operation
		for _, v := range apps {
			if v.Name == appName {
				LogMessage(fmt.Sprintf("Name: %v", v.Name), "g")
				LogMessage(fmt.Sprintf("Backend Port: %v", v.BackendPort), "g")
				LogMessage(fmt.Sprintf("Description: %v", v.Description), "g")
				LogMessage(fmt.Sprintf("Domain Name: %v", v.DomainName), "g")
				LogMessage(fmt.Sprintf("Protocol: %v", v.Protocol), "g")
				LogMessage(fmt.Sprintf("UUID: %v\n", v.UUID), "g")
				fmt.Println(v.Containers)
			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range apps {
		LogMessage(v.Name, "g")
	}
}

func showContainersAction(c *cli.Context) {
	api := getAPI(c)
	cntMeta, err := api.GetContainers()
	if err != nil {
		panic(err)
	}
	containers := cntMeta.Objects
	// check for container op
	containerID := c.String("id")
	if containerID != "" {
		// check for operation
		for _, v := range containers {
			if strings.Index(v.ContainerID, containerID) == 0 {
				LogMessage(fmt.Sprintf("ID: %v", v.ContainerID), "g")
				if v.Description != "" {
					LogMessage(fmt.Sprintf("Description: %v", v.Description), "g")
				}
				LogMessage(fmt.Sprintf("Image: %v", v.Meta.Config.Image), "g")
				LogMessage(fmt.Sprintf("CPU Shares: %v", v.Meta.Config.CpuShares), "g")
				LogMessage(fmt.Sprintf("Memory Limit: %v", v.Meta.Config.Memory), "g")
				LogMessage("Environment:", "g")
				LogMessage(fmt.Sprintf("  %v", strings.Join(v.Meta.Config.Env, "\n   ")), "")
				LogMessage(fmt.Sprintf("Created: %v", v.Meta.Created), "g")

			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range containers {
		LogMessage(fmt.Sprintf("%v (%v)", v.ContainerID, v.Meta.Config.Image), "g")
	}
}

func showHostsAction(c *cli.Context) {
	api := getAPI(c)
	meta, err := api.GetHosts()
	if err != nil {
		panic(err)
	}
	hosts := meta.Objects
	// check for container op
	name := c.String("name")
	if name != "" {
		// check for operation
		for _, v := range hosts {
			if v.Name == name {
				LogMessage(fmt.Sprintf("Name: %v", v.Name), "g")
				LogMessage(fmt.Sprintf("Hostname: %v", v.Hostname), "g")
				LogMessage(fmt.Sprintf("Port: %v", v.Port), "g")
				LogMessage(fmt.Sprintf("Enabled: %v", v.Enabled), "g")
			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range hosts {
		color := "g"
		if !v.Enabled {
			color = ""
		}
		LogMessage(fmt.Sprintf("%v (%v)", v.Name, v.Hostname), color)
	}
}

func configAction(c *cli.Context) {
	config, _ := loadConfig(c)
	LogMessage(fmt.Sprintf("Endpoint: %v", config.Url), "g")
	LogMessage(fmt.Sprintf("Username: %v", config.Username), "g")
	LogMessage(fmt.Sprintf("Version: %v", config.Version), "g")
	LogMessage(fmt.Sprintf("APIKey: %v...", config.ApiKey[0:5]), "g")
}
