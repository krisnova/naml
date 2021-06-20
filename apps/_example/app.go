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

package myapplication

import (
	"fmt"
	yamyams "github.com/kris-nova/yamyams/pkg"
	"k8s.io/client-go/kubernetes"
)

type MyApplication struct {
	meta      *yamyams.DeployableMeta
	resources []interface{}
}

func New() *MyApplication {
	return &MyApplication{
		meta: &yamyams.DeployableMeta{
			Name:        "Example Application",
			Command:     "example-app",
			Version:     "0.0.1",
			Description: "A simple example application",
		},
	}
}

func (v *MyApplication) Install(client *kubernetes.Clientset) error {
	return fmt.Errorf("[install] for %s not yet implemented", v.meta.Name)
	return nil
}

func (v *MyApplication) Uninstall(client *kubernetes.Clientset) error {
	return fmt.Errorf("[uninstall] for %s not yet implemented", v.meta.Name)
	return nil
}

func (v *MyApplication) Resources() []interface{} {
	return v.resources
}

func (v *MyApplication) About() *yamyams.DeployableMeta {
	return v.meta
}
