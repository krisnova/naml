//
// Copyright © 2022 Kris Nóva <kris@nivenly.com>
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

type MutatingWebhookConfiguration struct {
	KubeObject *admissionregistrationv1.MutatingWebhookConfiguration
	GoName     string
}

func NewMutatingWebhookConfiguration(obj *admissionregistrationv1.MutatingWebhookConfiguration) *MutatingWebhookConfiguration {
	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	return &MutatingWebhookConfiguration{
		KubeObject: obj,
		GoName:     goName(obj.Name),
	}
}

func (k MutatingWebhookConfiguration) Install() (string, []string) {
	c, err := Literal(k.KubeObject)
	if err != nil {
		logger.Debug(err.Error())
	}
	l := c.Source
	packages := c.Packages
	install := fmt.Sprintf(`
	{{ .GoName }}MutatingwebhookConfiguration := %s
	x.objects = append(x.objects, {{ .GoName }}MutatingwebhookConfiguration)

	if client != nil {
		_, err = client.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.TODO(), {{ .GoName }}MutatingwebhookConfiguration, v1.CreateOptions{})
		if err != nil {
			return err
		}
	}
`, l)

	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(install)
	buf := &bytes.Buffer{}
	k.KubeObject.Name = sanitizeK8sObjectName(k.KubeObject.Name)
	err = tpl.Execute(buf, k)
	if err != nil {
		logger.Debug(err.Error())
	}
	return alias(buf.String(), "admissionregistrationv1"), packages
}

func (k MutatingWebhookConfiguration) Uninstall() string {
	uninstall := `
	if client != nil {
		err = client.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.TODO(), "{{ .KubeObject.Name }}", metav1.DeleteOptions{})
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
