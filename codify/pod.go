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
	v1 "k8s.io/api/core/v1"
	"text/template"
)

type Pod struct {
	i *v1.Pod
}

func NewPod(pod *v1.Pod) *Pod {
	return &Pod{
		i: pod,
	}
}

func (p Pod) Install() string {
	install := `
	pod := &v1.Pod{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "{{ .Name }}",
			GenerateName:               "{{ .GenerateName }}",
			Namespace:                  "{{ .Namespace }}",
			ResourceVersion:            "{{ .ResourceVersion }}",
			Labels:                     map[string]string{},
			Annotations:                map[string]string{},
			ClusterName:                "{{ .ClusterName }}",
		},
		Spec:       v1.PodSpec{
			Volumes:                       nil,
			InitContainers:                nil,
			Containers:                    nil,
			EphemeralContainers:           nil,
			RestartPolicy:                 "{{ .Spec.RestartPolicy }}",
			TerminationGracePeriodSeconds: nil,
			ActiveDeadlineSeconds:         nil,
			DNSPolicy:                     "{{ .Spec.DNSPolicy }}",
			NodeSelector:                  nil,
			ServiceAccountName:            "{{ .Spec.ServiceAccountName }}",
			DeprecatedServiceAccount:      "{{ .Spec.DeprecatedServiceAccount }}",
			AutomountServiceAccountToken:  nil,
			NodeName:                      "{{ .Spec.NodeName }}",
			HostNetwork:                   {{ .Spec.HostNetwork }},
			HostPID:                       {{ .Spec.HostPID }},
			HostIPC:                       {{ .Spec.HostIPC }},
			ShareProcessNamespace:         nil,
			SecurityContext:               nil,
			ImagePullSecrets:              nil,
			Hostname:                      "{{ .Spec.Hostname }}",
			Subdomain:                     "{{ .Spec.Subdomain }}",
			Affinity:                      nil,
			SchedulerName:                 "{{ .Spec.SchedulerName }}",
			Tolerations:                   nil,
			HostAliases:                   nil,
			PriorityClassName:             "{{ .Spec.PriorityClassName }}",
			Priority:                      nil,
			DNSConfig:                     nil,
			ReadinessGates:                nil,
			RuntimeClassName:              nil,
			EnableServiceLinks:            nil,
			PreemptionPolicy:              nil,
			Overhead:                      nil,
			TopologySpreadConstraints:     nil,
			SetHostnameAsFQDN:             nil,
		},
	}
	_, err = client.CoreV1().Pods("{{ .Namespace }}").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
`
	tpl := template.New("pod")
	tpl.Parse(install)
	buf := &bytes.Buffer{}
	tpl.Execute(buf, p.i)
	return buf.String()
}

func (p Pod) Uninstall() string {
	uninstall := `
	err = client.CoreV1().Pods("{{ .Namespace }}").Delete(context.TODO(), "{{ .Name }}", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
 `
	tpl := template.New("dpod")
	tpl.Parse(uninstall)
	buf := &bytes.Buffer{}
	tpl.Execute(buf, p.i)
	return buf.String()
}
