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

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

	"github.com/kris-nova/logger"
)

type ValidatingwebhookConfiguration struct {
	KubeObject *admissionregistrationv1.ValidatingWebhookConfiguration
	GoName     string
}

func NewValidatingwebhookConfiguration(obj *admissionregistrationv1.ValidatingWebhookConfiguration) *ValidatingwebhookConfiguration {
	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	return &ValidatingwebhookConfiguration{
		KubeObject: obj,
		GoName:     goName(obj.Name),
	}
}

func (k ValidatingwebhookConfiguration) Install() (string, []string) {
	l, packages := Literal(k.KubeObject)
	install := fmt.Sprintf(`
	{{ .GoName }}ValidatingwebhookConfiguration := %s
	a.objects = append(a.objects, {{ .GoName }}ValidatingwebhookConfiguration)

	if client != nil {
		_, err = client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(context.TODO(), {{ .GoName }}ValidatingwebhookConfiguration, v1.CreateOptions{})
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
	return alias(buf.String(), "admissionregistrationv1"), packages
}

func (k ValidatingwebhookConfiguration) Uninstall() string {
	uninstall := `
	if client != nil {
		err = client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(context.TODO(), "{{ .KubeObject.Name }}", metav1.DeleteOptions{})
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
