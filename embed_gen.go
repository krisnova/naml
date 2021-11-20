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

const FormatMainGo string = `
// Copyright © {{ .CopyrightYear }} {{ .AuthorName }} {{ .AuthorEmail }}
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

package {{ .PackageName }}

import (
	"context"
	"fmt"
	"os"

{{ .Packages }}

	"github.com/kris-nova/naml"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// Version is the current release of your application.
var Version string = "{{ .Version }}"

func main() {
	// Load the application into the NAML registery
	// Note: naml.Register() can be used multiple times.
	naml.Register(NewApp("{{ .AppNameTitle }}", "{{ .Description }}"))

	// Run the generic naml command line program with
	// the application loaded.
	err := naml.RunCommandLine()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// App is a very important grown up business application.
type App struct {
	metav1.ObjectMeta
	description string
	objects []runtime.Object
	// ----------------------------------
	// Add your configuration fields here
	// ----------------------------------
}

// NewApp will create a new instance of App.
//
// See https://github.com/naml-examples for more examples.
//
// This is where you pass in fields to your application (similar to Values.yaml)
// Example: func NewApp(name string, example string, something int) *App
func NewApp(name, description string) *App {
	return &App{
		description: description,
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			ResourceVersion: Version,
		},
	    // ----------------------------------
	    // Add your configuration fields here
	    // ----------------------------------
	}
}

func (a *App) Install(client kubernetes.Interface) error {
	var err error
	{{ .Install }}
	return err
}

func (a *App) Uninstall(client kubernetes.Interface) error {
	var err error
	{{ .Uninstall }}
	return err
}

func (a *App) Description() string {
	return a.description
}

func (a *App) Meta() *metav1.ObjectMeta {
	return &a.ObjectMeta
}

func (a *App) Objects() []runtime.Object {
	return a.objects
}
`

const FormatLibraryGo string = `
// Copyright © {{ .CopyrightYear }} {{ .AuthorName }} {{ .AuthorEmail }}
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

package {{ .PackageName }}

import (
	"context"

{{ .Packages }}

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// {{ .AppNameTitle }}Version is the current release of your application.
var {{ .AppNameTitle }}Version string = "{{ .Version }}"

// App is a very important grown up business application.
type App struct {
	metav1.ObjectMeta
	description string
	objects []runtime.Object
	// ----------------------------------
	// Add your configuration fields here
	// ----------------------------------
}

// NewApp will create a new instance of App.
//
// See https://github.com/naml-examples for more examples.
//
// This is where you pass in fields to your application (similar to Values.yaml)
// Example: func NewApp(name string, example string, something int) *App
func NewApp(name, description string) *App {
	return &App{
		description: description,
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			ResourceVersion: {{ .AppNameTitle }}Version,
		},
	    // ----------------------------------
	    // Add your configuration fields here
	    // ----------------------------------
	}
}

func (a *App) Install(client kubernetes.Interface) error {
	var err error
	{{ .Install }}
	return err
}

func (a *App) Uninstall(client kubernetes.Interface) error {
	var err error
	{{ .Uninstall }}
	return err
}

func (a *App) Description() string {
	return a.description
}

func (a *App) Meta() *metav1.ObjectMeta {
	return &a.ObjectMeta
}

func (a *App) Objects() []runtime.Object {
	return a.objects
}
`

