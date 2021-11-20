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

package main

import (
	"context"
	"fmt"
	"os"

	{{ .Packages }}

	"github.com/kris-nova/naml"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

var Version string = "{{ .Version }}"

func main() {
	naml.Version = Version
	naml.Register(New{{ .AppNameTitle }}("{{ .AppNameTitle }}Instance", "{{ .Description }}"))
	err := naml.RunCommandLine()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type {{ .AppNameTitle }} struct {
	naml.AppMeta
	objects []runtime.Object
}

func New{{ .AppNameTitle }}(name, description string) *{{ .AppNameTitle }} {
	return &{{ .AppNameTitle }}{
		AppMeta: naml.AppMeta{
			Description: description,
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				ResourceVersion: Version,
			},
		},
	}
}

func (x *{{ .AppNameTitle }}) Install(client kubernetes.Interface) error {
	var err error
	{{ .Install }}
	return err
}

func (x *{{ .AppNameTitle }}) Uninstall(client kubernetes.Interface) error {
	var err error
	{{ .Uninstall }}
	return err
}

func (x *{{ .AppNameTitle }}) Meta() *naml.AppMeta {
	return &x.AppMeta
}

func (x *{{ .AppNameTitle }}) Objects() []runtime.Object {
	return x.objects
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

	"github.com/kris-nova/naml"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

var {{ .AppNameTitle }}Version string = "{{ .Version }}"

type {{ .AppNameTitle }} struct {
	naml.AppMeta
	objects []runtime.Object
}

func New{{ .AppNameTitle }}(name, description string) *{{ .AppNameTitle }} {
	return &{{ .AppNameTitle }}{
		AppMeta: naml.AppMeta{
			Description: description,
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				ResourceVersion: {{ .AppNameTitle }}Version,
			},
		},
	}
}

func (x *{{ .AppNameTitle }}) Install(client kubernetes.Interface) error {
	var err error
	{{ .Install }}
	return err
}

func (x *{{ .AppNameTitle }}) Uninstall(client kubernetes.Interface) error {
	var err error
	{{ .Uninstall }}
	return err
}

func (x *{{ .AppNameTitle }}) Meta() *naml.AppMeta {
	return &x.AppMeta
}

func (x *{{ .AppNameTitle }}) Objects() []runtime.Object {
	return x.objects
}
`

