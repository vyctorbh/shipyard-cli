package main

import (
    "github.com/codegangsta/cli"
    "fmt"
)

func applicationAction(c *cli.Context) {
    api := getAPI(c)
    appMeta, err := api.GetApplications()
    if err != nil { panic(err) }
    apps := appMeta.Objects
    // check for app op
    appName := c.String("name")
    if appName != "" {
        // check for operation
        if c.Bool("details") {
            for _, v := range apps {
                if v.Name == appName {
                    LogMessage(fmt.Sprintf("Name: %v", v.Name), "g")
                    LogMessage(fmt.Sprintf("Backend Port: %v", v.BackendPort), "g")
                    LogMessage(fmt.Sprintf("Description: %v", v.Description), "g")
                    LogMessage(fmt.Sprintf("Domain Name: %v", v.DomainName), "g")
                    LogMessage(fmt.Sprintf("Protocol: %v", v.Protocol), "g")
                    LogMessage(fmt.Sprintf("UUID: %v\n", v.UUID), "g")
                }
            }
            return
        }
    }
    // no op specified ; show all
    for _, v := range apps {
        LogMessage(v.Name, "g")
    }
}

func containerAction(c *cli.Context) {
    api := getAPI(c)
    cntMeta, err := api.GetContainers()
    if err != nil { panic(err) }
    containers := cntMeta.Objects
    // check for container op
    containerID := c.String("id")
    if containerID != "" {
        // check for operation
        if c.Bool("details") {
            for _, v := range containers {
                if v.ContainerID == containerID {
                    LogMessage(fmt.Sprintf("ID: %v", v.ContainerID), "g")
                    if v.Description != "" {
                        LogMessage(fmt.Sprintf("Description: %v", v.Description), "g")
                    }
                    LogMessage(fmt.Sprintf("Image: %v", v.Meta.Config.Image), "g")
                    LogMessage(fmt.Sprintf("Created: %v", v.Meta.Created), "g")

                }
            }
            return
        }
    }
    // no op specified ; show all
    for _, v := range containers {
        LogMessage(fmt.Sprintf("%v (%v)", v.ContainerID, v.Meta.Config.Image), "g")
    }
}
func configAction(c *cli.Context) {
    config, _ := loadConfig(c)
    LogMessage(fmt.Sprintf("Endpoint: %v", config.Url), "g")
    LogMessage(fmt.Sprintf("Username: %v", config.Username), "g")
    LogMessage(fmt.Sprintf("Version: %v", config.Version), "g")
    LogMessage(fmt.Sprintf("APIKey: %v...", config.ApiKey[0:5]), "g")
}
