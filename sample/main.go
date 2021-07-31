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

package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/naml"
	"k8s.io/client-go/kubernetes"
)

var Version string = "1.0.0"

func main() {
	naml.Register(NewApp("Barnaby", "A lovely sample application."))
	err := naml.RunCommandLine()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type App struct {
	metav1.ObjectMeta
	description string `json:"Description"`
	// --------------------
	// Add your fields here
	// --------------------
}

// NewApp will create a new instance of App.
//
// See https://github.com/naml-examples for more examples.
//
// Example: func NewApp(name string, example string, something int) *App
func NewApp(name, description string) *App {
	return &App{
		description: description,
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			ResourceVersion: Version,
		},
		// --------------------
		// Add your fields here
		// --------------------
	}
}

func (n *App) Install(client *kubernetes.Clientset) error {
	var err error

	var pod = &v1.Pod{
		TypeMeta: metav1.TypeMeta{
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
		Spec: v1.PodSpec{
			Volumes:        nil,
			InitContainers: nil,
			Containers: []v1.Container{
				v1.Container{
					Name:                     "",
					Image:                    "",
					Command:                  nil,
					Args:                     nil,
					WorkingDir:               "",
					Ports:                    nil,
					EnvFrom:                  nil,
					Env:                      nil,
					Resources:                v1.ResourceRequirements{},
					VolumeMounts:             nil,
					VolumeDevices:            nil,
					LivenessProbe:            nil,
					ReadinessProbe:           nil,
					StartupProbe:             nil,
					Lifecycle:                nil,
					TerminationMessagePath:   "",
					TerminationMessagePolicy: "",
					ImagePullPolicy:          "",
					SecurityContext:          nil,
					Stdin:                    false,
					StdinOnce:                false,
					TTY:                      false,
				},
			},
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

	return err
}

func (n *App) Uninstall(client *kubernetes.Clientset) error {
	var err error

	err = client.CoreV1().Pods("").Delete(context.TODO(), "", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return err
}

func (n *App) Description() string {
	return n.description
}

func (n *App) Meta() *metav1.ObjectMeta {
	return &n.ObjectMeta
}
