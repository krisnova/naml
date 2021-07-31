package naml

const FormatMainGo string = `
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

package main

import (
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/naml"
	"k8s.io/client-go/kubernetes"
)

var Version string = "1.0.0"

func main() {
	naml.Register(NewApp("Barnaby", "A lovely sample application."))
	err := naml.RunCommandLine()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type App struct {
	metav1.ObjectMeta
	description string
	// --------------------
	// Add your fields here
	// --------------------
}

// NewApp will create a new instance of App.
//
// See https://github.com/naml-examples for more examples.
//
// Example: func NewApp(name string, example string, something int) *App
func NewApp(name, description string) *App {
	return &App{
		description: description,
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			ResourceVersion: Version,
		},
		// --------------------
		// Add your fields here
		// --------------------
	}
}

func (n *App) Install(client *kubernetes.Clientset) error {
	var err error
	// [install]
	return err
}

func (n *App) Uninstall(client *kubernetes.Clientset) error {
	var err error
	// [uninstall]
	return err
}

func (n *App) Description() string {
	return n.description
}

func (n *App) Meta() *metav1.ObjectMeta {
	return &n.ObjectMeta
}
`
