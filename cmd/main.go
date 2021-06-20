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
	registry "github.com/kris-nova/yamyams"
	"github.com/kris-nova/yamyams/pkg"
	"os"

	"github.com/kris-nova/logger"
	"github.com/urfave/cli/v2"
)

var Version string = "---"

func main() {
	var verbose bool = true

	// cli assumes "-v" for version.
	// override that here
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	app := &cli.App{
		Name:      "YamYams",
		HelpName:  "yamyams",
		Usage:     "YAML alternative for managing Kubernetes packages directly with Go.",
		UsageText: " $ yamyams [options] <arguments>",
		Description: `YamYams gives infrastructure teams a framework to manage their applications directly with the Go programming language.
Instead of using YAML and templating we decided to just write Go. 
This is a command line tool that gives teams a starting point to start iterating.
Define applications in the /apps directory and register them in /registery.go.

Is there really that much of a difference with hard coding in Go versus writing YAML after all?`,
		Version: Version,
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
				UsageText:   "yamyams install [app]",
				Action: func(c *cli.Context) error {
					arguments := c.Args()
					if arguments.Len() != 1 {
						// Feature: We might want to have "yamyams install" just iterate through every application.
						cli.ShowCommandHelp(c, "install")
						List()
						os.Exit(1)
						return nil
					}
					appName := arguments.First()
					app := yamyams.Find(appName)
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
				UsageText:   "yamyams uninstall [app]",
				Action: func(c *cli.Context) error {
					arguments := c.Args()
					if arguments.Len() != 1 {
						// Feature: We might want to have "yamyams install" just iterate through every application.
						cli.ShowCommandHelp(c, "uninstall")
						List()
						os.Exit(1)
						return nil
					}
					appName := arguments.First()
					app := yamyams.Find(appName)
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
