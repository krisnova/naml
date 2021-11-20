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
	networkingv1 "k8s.io/api/networking/v1"
)

type IngressClass struct {
	KubeObject *networkingv1.IngressClass
	GoName     string
}

func NewIngressClass(obj *networkingv1.IngressClass) *IngressClass {
	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	return &IngressClass{
		KubeObject: obj,
		GoName:     goName(obj.Name),
	}
}

func (k IngressClass) Install() (string, []string) {
	l, packages := Literal(k.KubeObject)
	install := fmt.Sprintf(`
	{{ .GoName }}IngressClass := %s
	x.objects = append(x.objects, {{ .GoName }}IngressClass)

	if client != nil {
		_, err = client.NetworkingV1().IngressClasses().Create(context.TODO(), {{ .GoName }}IngressClass, v1.CreateOptions{})
		if err != nil {
			return err
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
	return alias(buf.String(), "networkingv1"), packages
}

func (k IngressClass) Uninstall() string {
	uninstall := `
	if client != nil {
		err = client.NetworkingV1().IngressClasses().Delete(context.TODO(), "{{ .KubeObject.Name }}", metav1.DeleteOptions{})
		if err != nil {
			return err
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
