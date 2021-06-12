//
// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, softwar
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
	"os"

	yamyam "github.com/kris-nova/yamyams/pkg"

	"github.com/kris-nova/logger"
	"k8s.io/client-go/util/homedir"
)

func main() {
	logger.BitwiseLevel = logger.LogEverything
	logger.Always("Hi hello yes welcome to my program.")
	logger.Always("This is a fun little deployment tool that lets you install yamyams in kubernetes.")
	logger.Always("Let's get started by making sure no bad yamyams are already installed.")
	logger.Always("This is called being idempotent!")
	logger.Always("Which is a fancy way of saying we always get the same outcome")
	logger.Always("Anyway... it basically just means whenever you run the code you get what you want.")
	y := yamyam.New()
	y.ContainerImage = "krisnova/yamyams"
	y.ContainerPort = 80
	err := y.KubernetesClient(fmt.Sprintf("%s/.kube/config", homedir.HomeDir()))
	if err != nil {
		// Oh no!
		logger.Critical("Unable to load the Kubernetes client")
		logger.Critical("Oof: %v", err)
		os.Exit(1) // <--- Kill the program
	}
	logger.Success("Success! Created Kubernetes Client!")
	err = y.UninstallKubernetes()
	if err != nil {
		// Oh no!
		logger.Warning("Not a huge deal, but just a heads up that something went wrong")
		logger.Warning("Looks like UninstallKubernetes() failed: %v", err)
	}
	logger.Success("Success! Cleaned up any existing deployments!")
	err = y.InstallKubernetes()
	if err != nil {
		// Oh no!
		logger.Critical("Unable to install in Kubernetes!")
		logger.Critical("Something went wrong: %v", err)
	}
	logger.Success("Success! Installed YamYams!")

}
