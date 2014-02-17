package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gcmurphy/getpass"
	"os"
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
		LogMessage("Error logging in.  Please check username/password.", "r")
		os.Exit(1)
	}
	saveConfig(username, userData.ApiKey, url, version)
	LogMessage("Login successful", "g")
}

func applicationsAction(c *cli.Context) {
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

func containersAction(c *cli.Context) {
	api := getAPI(c)
	showAll := c.Bool("all")
	cntMeta, err := api.GetContainers(showAll)
	if err != nil {
		panic(err)
	}
	containers := cntMeta.Objects
	// check for container op
	containerID := c.String("id")
	// check for operation
	restart := c.Bool("restart")
	stop := c.Bool("stop")
	start := c.Bool("start")
	remove := c.Bool("remove")
	// check for multiple requests
	if restart {
		if containerID == "" {
			LogMessage("Error: No container specified", "r")
			return
		}
		_, err := api.RestartContainer(containerID)
		if err != nil {
			LogMessage(fmt.Sprintf("Error restarting %s: %s", containerID, err), "r")
		} else {
			LogMessage(fmt.Sprintf("Restarted %s", containerID), "g")
		}
		return
	}
	if stop {
		if containerID == "" {
			LogMessage("Error: No container specified", "r")
			return
		}
		_, err := api.StopContainer(containerID)
		if err != nil {
			LogMessage(fmt.Sprintf("Error stopping %s: %s", containerID, err), "r")
		} else {
			LogMessage(fmt.Sprintf("Stopped %s", containerID), "g")
		}
		return
	}
	if start {
		if containerID == "" {
			LogMessage("Error: No container specified", "r")
			return
		}
		_, err := api.StartContainer(containerID)
		if err != nil {
			LogMessage(fmt.Sprintf("Error starting %s: %s", containerID, err), "r")
		} else {
			LogMessage(fmt.Sprintf("Started %s", containerID), "g")
		}
		return
	}
	if remove {
		if containerID == "" {
			LogMessage("Error: No container specified", "r")
			return
		}
		_, err := api.RemoveContainer(containerID)
		if err != nil {
			LogMessage(fmt.Sprintf("Error removing %s: %s", containerID, err), "r")
		} else {
			LogMessage(fmt.Sprintf("Removed %s", containerID), "g")
		}
		return
	}
	if containerID != "" {
		for _, v := range containers {
			if strings.Index(v.ContainerID, containerID) == 0 {
				LogMessage(fmt.Sprintf("ID: %v", v.ContainerID[:12]), "g")
				if v.Description != "" {
					LogMessage(fmt.Sprintf("Description: %v", v.Description), "g")
				}
				LogMessage(fmt.Sprintf("Image: %v", v.Meta.Config.Image), "g")
				LogMessage(fmt.Sprintf("CPU Shares: %v", v.Meta.Config.CpuShares), "g")
				LogMessage(fmt.Sprintf("Memory Limit: %v", v.Meta.Config.Memory), "g")
				if len(v.Meta.Config.Env) > 0 {
					LogMessage("Environment", "g")
					LogMessage(fmt.Sprintf("  %v", strings.Join(v.Meta.Config.Env, "\n   ")), "")
				}
				LogMessage(fmt.Sprintf("Created: %v", v.Meta.Created), "g")

			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range containers {
		LogMessage(fmt.Sprintf("%v %v", v.ContainerID[:12], v.Meta.Config.Image), "g")
	}
}

func imagesAction(c *cli.Context) {
	api := getAPI(c)
	meta, err := api.GetImages()
	if err != nil {
		panic(err)
	}
	images := meta.Objects
	// check for op
	id := c.String("id")
	if id != "" {
		for _, v := range images {
			if strings.Index(v.ID, id) == 0 {
				LogMessage(fmt.Sprintf("Repository: %v", v.Repository), "g")
			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range images {
		color := "g"
		LogMessage(fmt.Sprintf("%s %s", v.ID[:12], v.Repository), color)
	}
}

func hostsAction(c *cli.Context) {
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
	LogMessage(fmt.Sprintf("APIKey: %v", config.ApiKey), "g")
}
