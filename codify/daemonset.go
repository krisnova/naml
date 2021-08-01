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
	"text/template"
)

type DaemonSet struct {
	i *appsv1.DaemonSet
}

func NewDaemonSet(ds *appsv1.DaemonSet) *DaemonSet {
	ds.ObjectMeta = cleanObjectMeta(ds.ObjectMeta)
	ds.Status = appsv1.DaemonSetStatus{}
	return &DaemonSet{
		i: ds,
	}
}

func (k DaemonSet) Install() string {
	l := Literal(k.i)
	install := fmt.Sprintf(`
	{{ .Name }}DaemonSet := %s

	_, err = client.AppsV1().DaemonSets("{{ .Namespace }}").Create(context.TODO(), {{ .Name }}Deployment, v1.CreateOptions{})
	if err != nil {
		return err
	}
`, l)
	tpl := template.New("ds")
	tpl.Parse(install)
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, k.i)
	if err != nil {
		logger.Debug(err.Error())
	}
	return alias(buf.String(), "appsv1")
}

func (k DaemonSet) Uninstall() string {
	uninstall := `
	err = client.AppsV1().DaemonSets("{{ .Namespace }}").Delete(context.TODO(), "{{ .Name }}", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
 `
	tpl := template.New("dds")
	tpl.Parse(uninstall)
	buf := &bytes.Buffer{}
	k.i.Name = varName(k.i.Name)
	err := tpl.Execute(buf, k.i)
	if err != nil {
		logger.Debug(err.Error())
	}
	return buf.String()
}
