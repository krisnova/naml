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
	batchv1 "k8s.io/api/batch/v1"
)

type Job struct {
	KubeObject *batchv1.Job
	GoName     string
}

func NewJob(obj *batchv1.Job) *Job {
	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	obj.Status = batchv1.JobStatus{}
	return &Job{
		KubeObject: obj,
		GoName:     goName(obj.Name),
	}
}

func (k Job) Install() (string, []string) {
	l, packages := Literal(k.KubeObject)
	install := fmt.Sprintf(`
	{{ .GoName }}Job := %s
	a.objects = append(a.objects, {{ .GoName }}Job)

	if client != nil {
		_, err = client.BatchV1().Jobs("{{ .KubeObject.Namespace }}").Create(context.TODO(), {{ .GoName }}Job, v1.CreateOptions{})
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
	return alias(buf.String(), "batchv1"), packages
}

func (k Job) Uninstall() string {
	uninstall := `
	if client != nil {
		err = client.BatchV1().Jobs("{{ .KubeObject.Namespace }}").Delete(context.TODO(), "{{ .KubeObject.Name }}", metav1.DeleteOptions{})
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
