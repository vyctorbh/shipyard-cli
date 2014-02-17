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
	saveConfig(username, userData.ApiKey, url, version, c.GlobalString("config"))
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
	color := ""
	if appName != "" {
		// check for operation
		for _, v := range apps {
			if v.Name == appName {
				LogMessage(fmt.Sprintf("Name: %v", v.Name), color)
				LogMessage(fmt.Sprintf("Backend Port: %v", v.BackendPort), color)
				LogMessage(fmt.Sprintf("Description: %v", v.Description), color)
				LogMessage(fmt.Sprintf("Domain Name: %v", v.DomainName), color)
				LogMessage(fmt.Sprintf("Protocol: %v", v.Protocol), color)
				LogMessage(fmt.Sprintf("UUID: %v", v.UUID), color)
				if len(v.Containers) > 0 {
					LogMessage(fmt.Sprintf("Containers"), color)
					for _, c := range v.Containers {
						status := "running"
						if !c.Running {
							status = "stopped"
						}
						LogMessage(fmt.Sprintf("  %s %s (%s)", c.ContainerID[:12], c.Description, status), color)
					}
				}
			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range apps {
		LogMessage(v.Name, color)
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
	color := ""
	if containerID != "" {
		for _, v := range containers {
			if strings.Index(v.ContainerID, containerID) == 0 {
				status := "running"
				if !v.Running {
					status = "stopped"
				}
				LogMessage(fmt.Sprintf("ID: %s", v.ContainerID[:12]), color)
				LogMessage(fmt.Sprintf("Status: %s", status), color)
				LogMessage(fmt.Sprintf("Host: %v", v.Host.Name), color)
				if v.Description != "" {
					LogMessage(fmt.Sprintf("Description: %v", v.Description), color)
				}
				LogMessage(fmt.Sprintf("Image: %v", v.Meta.Config.Image), color)
				LogMessage(fmt.Sprintf("CPU Shares: %v", v.Meta.Config.CpuShares), color)
				LogMessage(fmt.Sprintf("Memory Limit: %v", v.Meta.Config.Memory), color)
				if len(v.Meta.Config.Env) > 0 {
					LogMessage("Environment", color)
					LogMessage(fmt.Sprintf("  %v", strings.Join(v.Meta.Config.Env, "\n   ")), color)
				}
				LogMessage(fmt.Sprintf("Created: %v", v.Meta.Created), color)

			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range containers {
		LogMessage(fmt.Sprintf("%s %s", v.ContainerID[:12], v.Meta.Config.Image), color)
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
	color := ""
	if id != "" {
		for _, v := range images {
			if strings.Index(v.ID, id) == 0 {
				LogMessage(fmt.Sprintf("ID: %s", v.ID), color)
				LogMessage(fmt.Sprintf("Host: %s", v.Host.Name), color)
				LogMessage(fmt.Sprintf("Repository: %s", v.Repository), color)
			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range images {
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
	color := ""
	if name != "" {
		// check for operation
		for _, v := range hosts {
			if v.Name == name {
				LogMessage(fmt.Sprintf("Name: %v", v.Name), color)
				LogMessage(fmt.Sprintf("Hostname: %v", v.Hostname), color)
				LogMessage(fmt.Sprintf("Port: %v", v.Port), color)
				LogMessage(fmt.Sprintf("Enabled: %v", v.Enabled), color)
			}
		}
		return
	}
	// no op specified ; show all
	for _, v := range hosts {
		if !v.Enabled {
			color = ""
		}
		LogMessage(fmt.Sprintf("%v (%v)", v.Name, v.Hostname), color)
	}
}

func configAction(c *cli.Context) {
	config, _ := loadConfig(c)
	color := ""
	LogMessage(fmt.Sprintf("Endpoint: %v", config.Url), color)
	LogMessage(fmt.Sprintf("Username: %v", config.Username), color)
	LogMessage(fmt.Sprintf("Version: %v", config.Version), color)
	LogMessage(fmt.Sprintf("APIKey: %v", config.ApiKey), color)
}
