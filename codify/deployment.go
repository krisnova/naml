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
	"github.com/kris-nova/logger"
	appsv1 "k8s.io/api/apps/v1"
	"strings"
	"text/template"
	"time"
)

type Deployment struct {
	KubeObject *appsv1.Deployment
	GoName string
}

func NewDeployment(obj *appsv1.Deployment) *Deployment {
	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	obj.Status = appsv1.DeploymentStatus{}
	return &Deployment{
		KubeObject: obj,
		GoName:     strings.ReplaceAll(obj.Name, "-", "_"),
	}
}

func (k Deployment) Install() string {
	l := Literal(k.KubeObject)
	install := fmt.Sprintf(`
	// Adding a deployment: "{{ .KubeObject.Name }}"
	{{ .GoName }}Deployment := %s

	_, err = client.AppsV1().Deployments("{{ .KubeObject.Namespace }}").Create(context.TODO(), {{ .GoName }}Deployment, v1.CreateOptions{})
	if err != nil {
		return err
	}
`, l)
	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(install)
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, k)
	if err != nil {
		logger.Debug(err.Error())
	}
	return alias(buf.String(), "appsv1")
}

func (k Deployment) Uninstall() string {
	uninstall := `
	err = client.AppsV1().Deployments("{{ .KubeObject.Namespace }}").Delete(context.TODO(), "{{ .KubeObject.Name }}", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
 `
	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(uninstall)
	buf := &bytes.Buffer{}
	k.KubeObject.Name = sanitizeK8sObjectName(k.KubeObject.Name)
	err := tpl.Execute(buf, k)
	if err != nil {
		logger.Debug(err.Error())
	}
	return buf.String()
}