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
	"path"
	"strings"

	"k8s.io/client-go/util/homedir"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	KubeconfigEnvironmentalVariable = "KUBECONFIG"
	KubeconfigDefaultDirectory      = ".kube"
	KubeconfigDefaultFile           = "config"
)

// cachedClient is package level state that will cache the Clientset
var cachedClient *kubernetes.Clientset

// kubeConfigPathValue is a very flexible string, so we have to handle it
// carefully. It can be either the environmental variable, or the CLI flag
//
// KUBECONFIG environmental variable
//
// This can be a list of configs such as:
//     config:config-sample:config-boops
//
// --kubeconfig
// This can be a single path such as
//     ~/.kube/config
//     /home/me/.kube/config
//     config (which should assume ${HOME}.kube/config)
var kubeConfigPathValue string

// Client is used to authenticate with Kubernetes and build the Kube client
// for the rest of the program.
func Client() (*kubernetes.Clientset, error) {
	if cachedClient != nil {
		return cachedClient, nil
	}

	// Windows support requires us to split on ";" instead of ":"
	// I personally have no use case to care for this. Pull requests accepted.
	configs := strings.Split(kubeConfigPathValue, ":")
	if len(configs) > 1 {
		for _, config := range configs {
			client, err := ClientFromPath(path.Join(homedir.HomeDir(), KubeconfigDefaultDirectory, config))
			if err == nil {
				// Just pick the first one
				return client, nil
			}
		}
		// If we get here, we have no valid config
		// We can just silently try and fail...
	}
	client, err := ClientFromPath(kubeConfigPathValue)
	if err != nil {
		return nil, fmt.Errorf("unable to load kube config: %v", err)
	}
	cachedClient = client
	return cachedClient, nil
}

// ClientFromPath is used to authenticate with Kubernetes and build the Kube client
// for the rest of the program given a specific kube config path.
//
// Useful for testing.
func ClientFromPath(kubeConfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to find local kube config [%s]: %v", kubeConfigPath, err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to build kube config: %v", err)
	}
	return client, nil
}

// ClientFromFlags will plumb well-known command line flags through to the kubeconfig
func ClientFromFlags(apiUrl, kubeConfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiUrl, kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to find local kube config [%s]: %v", kubeConfigPath, err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to build kube config: %v", err)
	}
	return client, nil
}
