//
// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//    ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗
//    ████╗  ██║██╔═████╗██║   ██║██╔══██╗
//    ██╔██╗ ██║██║██╔██║██║   ██║███████║
//    ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║
//    ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║
//    ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝

package main

import (
	"fmt"
	registry "github.com/kris-nova/naml"
	"github.com/kris-nova/naml/pkg"
	"os"

	"github.com/kris-nova/logger"
	"github.com/urfave/cli/v2"
)

func main() {
	var verbose bool = true

	// cli assumes "-v" for version.
	// override that here
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Print the version",
	}

	app := &cli.App{
		Name:      "naml",
		HelpName:  "naml",
		Usage:     "YAML alternative for managing Kubernetes packages directly with Go.",
		UsageText: " $ naml [options] <arguments>",
		Description: `
Use naml to start encapsulating your applications with Go.
Take advantage of all the lovely features of the Go programming language.

Is there really THAT much of a difference with defining an application in Go compared to defining an application in YAML after all?`,
		Version: registry.Version,
		Authors: []*cli.Author{
			{
				Name:  "Kris Nóva",
				Email: "kris@nivenly.com",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Value:       true,
				Usage:       "toggle verbose mode for logger.",
				Destination: &verbose,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "install",
				Aliases:     []string{"i"},
				Description: "Will execute the Install method for a specific app.",
				Usage:       "Install a package in Kubernetes.",
				UsageText:   "naml install [app]",
				Action: func(c *cli.Context) error {
					arguments := c.Args()
					if arguments.Len() != 1 {
						// Feature: We might want to have "naml install" just iterate through every application.
						cli.ShowCommandHelp(c, "install")
						List()
						os.Exit(1)
						return nil
					}
					appName := arguments.First()
					app := naml.Find(appName)
					if app == nil {
						return fmt.Errorf("Invalid application name (Application not registered): %s", appName)
					}
					logger.Info("Installing [%s]", appName)
					return Install(app)
				},
			},
			{
				Name:        "uninstall",
				Aliases:     []string{"u"},
				Description: "Will execute the Uninstall method for a specific app.",
				Usage:       "Uninstall a package in Kubernetes",
				UsageText:   "naml uninstall [app]",
				Action: func(c *cli.Context) error {
					arguments := c.Args()
					if arguments.Len() != 1 {
						// Feature: We might want to have "naml install" just iterate through every application.
						cli.ShowCommandHelp(c, "uninstall")
						List()
						os.Exit(1)
						return nil
					}
					appName := arguments.First()
					app := naml.Find(appName)
					if app == nil {
						return fmt.Errorf("Invalid application name (Application not registered): %s", appName)
					}
					logger.Info("Uninstalling [%s]", appName)
					return Uninstall(app)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "[local] List applications.",
				Action: func(c *cli.Context) error {
					List()
					return nil
				},
			},
		},
	}

	if verbose {
		logger.BitwiseLevel = logger.LogEverything
		logger.Always("[Verbose Mode]")
	} else {
		logger.BitwiseLevel = logger.LogAlways | logger.LogCritical | logger.LogWarning | logger.LogDeprecated
	}

	// Load whatever apps are defined in registry.go
	registry.Load()

	err := app.Run(os.Args)
	if err != nil {
		logger.Critical("%v", err)
	}
	os.Exit(0)
}
