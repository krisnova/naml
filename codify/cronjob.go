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
	batchv1 "k8s.io/api/batch/v1"
	"text/template"
	"time"
)

type CronJob struct {
	KubeObject *batchv1.CronJob
}

func NewCronJob(obj *batchv1.CronJob) *CronJob {
	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	obj.Status = batchv1.CronJobStatus{}
	return &CronJob{
		KubeObject: obj,
	}
}

func (k CronJob) Install() string {
	l := Literal(k.KubeObject)
	install := fmt.Sprintf(`
	{{ .Name }}CronJob := %s

	_, err = client.BatchV1().CronJobs("{{ .Namespace }}").Create(context.TODO(), {{ .Name }}CronJob, v1.CreateOptions{})
	if err != nil {
		return err
	}
`, l)

	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(install)
	buf := &bytes.Buffer{}
	k.KubeObject.Name = sanitizeK8sObjectName(k.KubeObject.Name)
	err := tpl.Execute(buf, k.KubeObject)
	if err != nil {
		logger.Debug(err.Error())
	}
	return alias(buf.String(), "batchv1")
}

func (k CronJob) Uninstall() string {
	uninstall := `
	err = client.BatchV1().CronJobs("{{ .Namespace }}").Delete(context.TODO(), "{{ .Name }}", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
 `
	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(uninstall)
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, k.KubeObject)
	if err != nil {
		logger.Debug(err.Error())
	}
	return buf.String()
}
