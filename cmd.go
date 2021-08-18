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
//   ███╗   ██╗ █████╗ ███╗   ███╗██╗
//   ████╗  ██║██╔══██╗████╗ ████║██║
//   ██╔██╗ ██║███████║██╔████╔██║██║
//   ██║╚██╗██║██╔══██║██║╚██╔╝██║██║
//   ██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗
//   ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝
//

package naml

import (
	"fmt"
	"os"
	"strings"
	"time"

	"k8s.io/client-go/util/homedir"

	"github.com/kris-nova/logger"
	"github.com/urfave/cli/v2"
)

func RunCommandLineAndExit() {
	err := RunCommandLine()
	if err != nil {
		logger.Critical(err.Error())
		os.Exit(1)
	}
}

// RunCommandLine is the global NAML command line program.
//
// Use this if you would like to use the built in NAML command line interface.
func RunCommandLine() error {
	// Default options
	return RunCommandLineWithOptions()
}

// RunCommandLineWithOptions is here so we can default values in RunCommandLine() that
// we would want to pass in here later (tests, etc)
func RunCommandLineWithOptions() error {
	// with is a set of paths that the user has specificed for naml
	// to run with
	var with cli.StringSlice

	// verbose is the logger verbosity
	var verbose bool = false

	// kubeconfig is the --kubeconfig value
	// which is used in our Client() code
	var kubeconfig string

	// codifyAppNameRaw is the app name passed in the raw form
	var codifyAppNameRaw string

	// output will be what output type for the output subcommand
	var output string

	codifyValues := &MainGoValues{
		AuthorEmail:   "<kris@nivenly.com>",
		AuthorName:    "Kris Nóva",
		CopyrightYear: fmt.Sprintf("%d", time.Now().Year()),
		AppNameLower:  "app",
		AppNameTitle:  "App",
		Version:       "0.0.1",
		Description:   "very serious grown up business application does important beep boops",
	}

	// cli assumes "-v" for version.
	// override that here
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Print the version",
	}

	// ********************************************************
	// [ NAML APPLICATION ]
	// ********************************************************

	app := &cli.App{
		Name:     "naml",
		HelpName: "Not Another Markup Language",
		Usage:    "Framework for managing Kubernetes applications in Go.",
		UsageText: `Use naml like any command line tool.
		naml [options] command [arguments...]

	Use naml to convert Kubernetes YAML to Go.
		kubectl get deployment -o yaml | naml codify

	Use naml to embed applications into a binary.
		naml list

	Use naml to install and uninstall in Kubernetes.
		naml install   <app>
		naml uninstall <app>
`,
		Action: func(context *cli.Context) error {
			if output != "" {
				return outputFunc(output, context)
			}
			Banner()
			cli.ShowSubcommandHelp(context)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Value:       false,
				Usage:       "toggle verbose mode for logger",
				Destination: &verbose,
			},
			&cli.StringFlag{
				Name:        "kubeconfig",
				Value:       "~/.kube/config",
				Usage:       "Kubeconfig path (default: ~/.kube/config)",
				Destination: &kubeconfig,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "Output encoded. (yaml, json)",
				Destination: &output,
			},
			&cli.StringSliceFlag{
				Name:        "with",
				Aliases:     []string{"f", "w"}, // use -f to follow kubectl -f syntax trolol
				Usage:       "include other naml binaries",
				Destination: &with,
			},
		},
		Commands: []*cli.Command{

			// ********************************************************
			// [ INSTALL ]
			// ********************************************************

			{
				Name:      "install",
				Aliases:   []string{"i"},
				Usage:     "Install a package in Kubernetes",
				UsageText: "naml install [name]",
				Action: func(c *cli.Context) error {
					// ----------------------------------
					err := AllInit(kubeconfig, verbose, with.Value())
					if err != nil {
						return err
					}
					// ----------------------------------

					// Right away if it's just one app use it
					if len(Registry()) == 1 {
						for _, app := range Registry() {
							return Install(app)
						}
					}

					arguments := c.Args()
					if arguments.Len() != 1 {
						Banner()
						cli.ShowCommandHelp(c, "install")
						List()
						return nil
					}
					appName := arguments.First()
					app := Find(appName)
					if app == nil {
						return fmt.Errorf("Invalid application name (Application not registered): %s", appName)
					}
					logger.Info("Installing [%s]", appName)
					return Install(app)
				},
			},

			// ********************************************************
			// [ UNINSTALL ]
			// ********************************************************

			{
				Name:      "uninstall",
				Aliases:   []string{"u"},
				Usage:     "Uninstall a package in Kubernetes",
				UsageText: "naml uninstall [name]",
				Action: func(c *cli.Context) error {
					// ----------------------------------
					err := AllInit(kubeconfig, verbose, with.Value())
					if err != nil {
						return err
					}
					// ----------------------------------

					// Right away if it's just one app use it
					if len(Registry()) == 1 {
						for _, app := range Registry() {
							return Uninstall(app)
						}
					}

					arguments := c.Args()
					if arguments.Len() != 1 {
						// Feature: We might want to have "naml install" just iterate through every application.
						Banner()
						cli.ShowCommandHelp(c, "uninstall")
						List()
						return nil
					}
					appName := arguments.First()
					app := Find(appName)
					if app == nil {
						return fmt.Errorf("Invalid application name (Application not registered): %s", appName)
					}
					logger.Info("Uninstalling [%s]", appName)
					return Uninstall(app)
				},
			},

			// ********************************************************
			// [ CODIFY ]
			// ********************************************************

			{
				Name:      "codify",
				Aliases:   []string{"c"},
				Usage:     "Convert YAML to syntactically correct Go code",
				UsageText: "cat deploy.yaml | naml codify > main.go",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "author-name",
						Value:       "Kris Nóva",
						Usage:       "Name for the copyright header",
						Destination: &codifyValues.AuthorName,
					},
					&cli.StringFlag{
						Name:        "author-email",
						Value:       "<kris@nivenly.com>",
						Usage:       "Email for the copyright header",
						Destination: &codifyValues.AuthorEmail,
					},
					&cli.StringFlag{
						Name:        "description",
						Value:       fmt.Sprintf("Application autogenerated from NAML v%s", Version),
						Usage:       "Description for the application",
						Destination: &codifyValues.Description,
					},
					&cli.StringFlag{
						Name:        "name",
						Value:       "App",
						Usage:       "Name to use for the application",
						Destination: &codifyAppNameRaw,
					},
				},
				Action: func(c *cli.Context) error {
					// ----------------------------------
					err := AllInit(kubeconfig, verbose, with.Value())
					if err != nil {
						return err
					}
					// ----------------------------------

					codifyValues.AppNameLower = strings.ToLower(codifyAppNameRaw)
					codifyValues.AppNameTitle = strings.Title(codifyValues.AppNameLower)

					cbytes, err := Codify(os.Stdin, codifyValues)
					if err != nil {
						// Codify prints to stderr
						fmt.Fprintf(os.Stderr, "Error during codify: %v", err)
						return err
					}
					fmt.Println(string(cbytes))
					return nil
				},
			},

			// ********************************************************
			// [ LIST ]
			// ********************************************************

			{
				Name:      "list",
				Aliases:   []string{"l"},
				Usage:     "List embedded applications",
				UsageText: "naml list",
				Action: func(c *cli.Context) error {
					// ----------------------------------
					err := AllInit(kubeconfig, verbose, with.Value())
					if err != nil {
						return err
					}
					// ----------------------------------
					Banner()
					List()
					return nil
				},
			},

			// ********************************************************
			// [ OUTPUT ]
			// ********************************************************

			{
				Name:      "output",
				Aliases:   []string{"o"},
				Usage:     "Output embedded applications. (yaml, json)",
				UsageText: "naml output [name] -o yaml",
				Action: func(c *cli.Context) error {
					return outputFunc(output, c)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Value:       "yaml",
						Usage:       "Output encoding. (yaml, json)",
						Destination: &output,
					},
				},
			},

			// ********************************************************
			// [ RPC ]
			// ********************************************************

			{
				Name:      "rpc",
				Aliases:   []string{"r"},
				Usage:     "Run in RPC mode. Only run this if you know what you are doing",
				UsageText: "naml rpc",
				Action: func(c *cli.Context) error {
					err := RunRPC()
					if err != nil {
						return fmt.Errorf("unable to run in runtime mode: %v", err)
					}
					return nil
				},
			},
		},
	}
	return app.Run(os.Args)
}

func outputFunc(encoding string, c *cli.Context) error {
	var o OutputEncoding
	encoding = strings.ToLower(encoding)
	if encoding == "json" {
		o = OutputJSON
	}
	if encoding == "yaml" {
		o = OutputYAML
	}

	arguments := c.Args()

	// No specific apps were passed
	if arguments.Len() != 1 {
		for name, _ := range Registry() {
			err := RunOutput(name, o)
			if err != nil {
				return fmt.Errorf("unable to run in runtime mode: %v", err)
			}
		}
		return nil
	}
	appName := arguments.First()
	err := RunOutput(appName, o)
	if err != nil {
		return fmt.Errorf("unable to run in runtime mode: %v", err)
	}
	return nil
}

// AllInit is the "constructor" for every command line flag.
// This is how we use naml -w to include sub-namls
func AllInit(kubeConfigPath string, verbose bool, with []string) error {

	// [ Verbosity System ]
	if verbose {
		logger.BitwiseLevel = logger.LogEverything
		logger.Always("*** [ Verbose Mode ] ***")
	} else {
		logger.BitwiseLevel = logger.LogAlways | logger.LogCritical | logger.LogWarning | logger.LogDeprecated
	}

	// [ Kubeconfig System ]
	// 1. Check if environmental variable is set
	// 2. Default to the --kubeconfig flag
	// 3. Follow the logic in the Clientcmd (path, masterURL, inCluster, default)

	// Format "~" in command line string
	kubeConfigPath = strings.ReplaceAll(kubeConfigPath, "~", homedir.HomeDir())

	// Here be dragons
	// We probably need an entire fucking client package, but for now
	// this will get us to 1.0.0
	envVarValue := os.Getenv(KubeconfigEnvironmentalVariable)
	if envVarValue == "" {
		kubeConfigPathValue = kubeConfigPath
	} else {
		kubeConfigPathValue = envVarValue
	}
	logger.Debug("Kubeconfig Value: %s", kubeConfigPathValue)

	// [ Child Runtime System ]
	if len(with) > 0 {
		for _, childPath := range with {
			for i := 0; i < 3; i++ {
				err := AddRPC(childPath)
				if err != nil {
					logger.Warning("Unable to add child naml %s: %v", childPath, err)
					time.Sleep(time.Millisecond * 20)
				} else {
					break
				}
			}

		}
	}

	// If running naml with children, register them with the registry
	if len(remotes) > 0 {
		err := RegisterRemoteApplications()
		if err != nil {
			return fmt.Errorf("unable to register children: %v", err)
		}
	}
	return nil
}

// Install is used to install an application in Kubernetes
func Install(app Deployable) error {

	// Check if app is a remote app
	if remoteApp, ok := app.(*RPCApplication); ok {
		// Do NOT pass in this local kubernetes client!
		return remoteApp.Install(nil)
	}

	// Only grab a client if we are running in this instance!
	client, err := Client()
	if err != nil {
		return err
	}

	// Install
	err = app.Install(client)
	if err != nil {
		return err
	}
	meta := app.Meta()
	if meta.Namespace == "" {
		meta.Namespace = "default"
	}
	fmt.Printf("Installed %s in namespace [%s]\n", meta.Name, meta.Namespace)
	logger.Success("Successfully installed [%s]", app.Meta().Name)
	return nil
}

// List the naml package information in stdout
func List() {
	fmt.Println("")
	for _, app := range Registry() {
		fmt.Printf("[%s]\n", app.Meta().Name)
		fmt.Printf("  description : %s\n", app.Description())
		fmt.Printf("  version     : %s\n", app.Meta().ResourceVersion)
		fmt.Printf("\n")
	}
}

// Uninstall is used to uninstall an application in Kubernetes
func Uninstall(app Deployable) error {

	// Check if app is a remote app
	if remoteApp, ok := app.(*RPCApplication); ok {
		// Do NOT pass in this local kubernetes client!
		return remoteApp.Uninstall(nil)
	}

	// Only grab a client if we are running in this instance!

	client, err := Client()
	if err != nil {
		return err
	}

	// Uninstall
	err = app.Uninstall(client)
	if err != nil {
		return err
	}
	meta := app.Meta()
	if meta.Namespace == "" {
		meta.Namespace = "default"
	}
	fmt.Printf("Uninstalled %s in namespace [%s]\n", meta.Name, meta.Namespace)
	return nil
}
