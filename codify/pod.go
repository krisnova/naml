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

import v1 "k8s.io/api/core/v1"

type Pod struct {

}

func NewPod(pod *v1.Pod) *Pod {
	return &Pod{}
}

func (p Pod) Install() string {
	return `

	pod := &v1.Pod{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "",
			GenerateName:               "",
			Namespace:                  "",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     map[string]string{},
			Annotations:                map[string]string{},
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ClusterName:                "",
			ManagedFields:              nil,
		},
		Spec:       v1.PodSpec{
			Volumes:                       nil,
			InitContainers:                nil,
			Containers:                    nil,
			EphemeralContainers:           nil,
			RestartPolicy:                 "",
			TerminationGracePeriodSeconds: nil,
			ActiveDeadlineSeconds:         nil,
			DNSPolicy:                     "",
			NodeSelector:                  nil,
			ServiceAccountName:            "",
			DeprecatedServiceAccount:      "",
			AutomountServiceAccountToken:  nil,
			NodeName:                      "",
			HostNetwork:                   false,
			HostPID:                       false,
			HostIPC:                       false,
			ShareProcessNamespace:         nil,
			SecurityContext:               nil,
			ImagePullSecrets:              nil,
			Hostname:                      "",
			Subdomain:                     "",
			Affinity:                      nil,
			SchedulerName:                 "",
			Tolerations:                   nil,
			HostAliases:                   nil,
			PriorityClassName:             "",
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
	_, err = client.CoreV1().Pods("").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

`
}

func (p Pod) Uninstall() string {
	return `
	err = client.CoreV1().Pods("").Delete(context.TODO(), "", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
 `
}
