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

package yamyams

import (
	"fmt"
	"k8s.io/client-go/util/homedir"
	"path"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

const (
	TestClusterName string = "yamyamstestcluster"
)

var (
	isStarted      bool   = false
	kubeConfigPath string = path.Join(homedir.HomeDir(), ".kube", "yamyams.conf")
)

func TestClusterStart() error {
	if isStarted {
		return nil
	}
	provider := cluster.NewProvider(cluster.ProviderWithDocker(), cluster.ProviderWithLogger(cmd.NewLogger()))
	err := provider.Create(TestClusterName)
	if err != nil {
		return fmt.Errorf("unable to create kind test cluster: %v", err)
	}
	err = provider.ExportKubeConfig(TestClusterName, kubeConfigPath)
	if err != nil {
		return fmt.Errorf("unable to export test kube config: %v", err)
	}
	isStarted = true
	return nil
}

func TestClusterKubeConfigPath() string {
	return kubeConfigPath
}

func TestClusterStop() error {
	provider := cluster.NewProvider(cluster.ProviderWithDocker())
	err := provider.Delete(TestClusterName, kubeConfigPath)
	if err != nil {
		return fmt.Errorf("unable to delete kind test cluster: %v", err)
	}
	return nil
}
