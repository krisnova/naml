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

package codify

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/kris-nova/logger"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type CustomResourceDefinition struct {
	KubeObject *apiextensionsv1.CustomResourceDefinition
	GoName     string
}

func NewCustomResourceDefinition(obj *apiextensionsv1.CustomResourceDefinition) *CustomResourceDefinition {

	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	if obj.Namespace == "" {
		obj.Namespace = "default"
	}
	return &CustomResourceDefinition{
		KubeObject: obj,
		GoName:     goName(obj.Name),
	}
}

func (k CustomResourceDefinition) Install() (string, []string) {
	l, packages := Literal(k.KubeObject)
	install := fmt.Sprintf(`
	{{ .GoName }}CustomResourceDefinition := %s
	a.objects = append(a.objects, {{ .GoName }}CustomResourceDefinition)
	
	if client != nil {
		result := client.ExtensionsV1beta1().RESTClient().Post().Namespace("{{ .KubeObject.Namespace }}").Body({{ .GoName }}CustomResourceDefinition).Do(context.TODO())
		if result.Error() != nil {
			return result.Error()
		}
	}
`, l)

	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(install)
	buf := &bytes.Buffer{}
	k.KubeObject.Name = sanitizeK8sObjectName(k.KubeObject.Name)
	err := tpl.Execute(buf, k)
	if err != nil {
		logger.Debug(err.Error())
	}
	return alias(buf.String(), "apiextensionsv1"), packages
}

func (k CustomResourceDefinition) Uninstall() string {
	uninstall := `
	if client != nil {
		result := client.ExtensionsV1beta1().RESTClient().Delete().Namespace(a.Namespace).Name(a.Name).Do(context.TODO())
		if result.Error() != nil {
			return result.Error()
		}
	}
 `

	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(uninstall)
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, k)
	if err != nil {
		logger.Debug(err.Error())
	}
	return buf.String()
}
