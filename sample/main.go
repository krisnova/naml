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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"

	"github.com/kris-nova/naml"
	"k8s.io/client-go/kubernetes"
)

var Version string = "0.0.1"

func main() {
	// Load the application into the NAML registery
	// Note: naml.Register() can be used multiple times.
	//
	naml.Register(NewApp("App", "very serious grown up business application does important beep boops"))

	// Run the generic naml command line program with
	// the application loaded.
	err := naml.RunCommandLine()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type App struct {
	metav1.ObjectMeta
	description string
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

func (a *App) Install(client *kubernetes.Clientset) error {
	var err error

	boopsDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/appsv1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:         "boops",
			GenerateName: "",
			Namespace:    "default",
			Labels:       map[string]string{"app": "boops"},
			Annotations:  map[string]string{"deployment.kubernetes.io/revision": "1"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: naml.I32p(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "",
					Labels: map[string]string{"app": "boops"},
				},
				Spec: corev1.PodSpec{
					Volumes:        []corev1.Volume(nil),
					InitContainers: []corev1.Container(nil),
					Containers: []corev1.Container{
						{
							Name:       "nginx",
							Image:      "nginx",
							Command:    []string(nil),
							Args:       []string(nil),
							WorkingDir: "",
							Ports:      []corev1.ContainerPort(nil),
							EnvFrom:    []corev1.EnvFromSource(nil),
							Env:        []corev1.EnvVar(nil),
							Resources: corev1.ResourceRequirements{Limits: corev1.ResourceList(nil),
								Requests: corev1.ResourceList(nil)},
							VolumeMounts:             []corev1.VolumeMount(nil),
							VolumeDevices:            []corev1.VolumeDevice(nil),
							LivenessProbe:            (*corev1.Probe)(nil),
							ReadinessProbe:           (*corev1.Probe)(nil),
							StartupProbe:             (*corev1.Probe)(nil),
							Lifecycle:                (*corev1.Lifecycle)(nil),
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: "File",
							ImagePullPolicy:          "Always",
							SecurityContext:          (*corev1.SecurityContext)(nil),
							Stdin:                    false,
							StdinOnce:                false,
							TTY:                      false}},
					EphemeralContainers:       []corev1.EphemeralContainer(nil),
					RestartPolicy:             "Always",
					ActiveDeadlineSeconds:     (*int64)(nil),
					DNSPolicy:                 "ClusterFirst",
					NodeSelector:              map[string]string(nil),
					ServiceAccountName:        "",
					NodeName:                  "",
					HostNetwork:               false,
					HostPID:                   false,
					HostIPC:                   false,
					ImagePullSecrets:          []corev1.LocalObjectReference(nil),
					Hostname:                  "",
					Subdomain:                 "",
					Affinity:                  (*corev1.Affinity)(nil),
					SchedulerName:             "default-scheduler",
					Tolerations:               []corev1.Toleration(nil),
					HostAliases:               []corev1.HostAlias(nil),
					PriorityClassName:         "",
					Priority:                  (*int32)(nil),
					DNSConfig:                 (*corev1.PodDNSConfig)(nil),
					ReadinessGates:            []corev1.PodReadinessGate(nil),
					RuntimeClassName:          (*string)(nil),
					EnableServiceLinks:        (*bool)(nil),
					PreemptionPolicy:          (*corev1.PreemptionPolicy)(nil),
					Overhead:                  corev1.ResourceList(nil),
					TopologySpreadConstraints: []corev1.TopologySpreadConstraint(nil),
					SetHostnameAsFQDN:         (*bool)(nil)}},
			Strategy: appsv1.DeploymentStrategy{
				Type: "RollingUpdate",
			},
			MinReadySeconds: 0,
			Paused:          false,
		},
		Status: appsv1.DeploymentStatus{ObservedGeneration: 0,
			Replicas:            0,
			UpdatedReplicas:     0,
			ReadyReplicas:       0,
			AvailableReplicas:   0,
			UnavailableReplicas: 0,
			Conditions:          []appsv1.DeploymentCondition(nil),
		},
	}

	_, err = client.AppsV1().Deployments("default").Create(context.TODO(), boopsDeployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return err
}

func (a *App) Uninstall(client *kubernetes.Clientset) error {
	var err error

	err = client.AppsV1().Deployments("default").Delete(context.TODO(), "boops", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return err
}

func (a *App) Description() string {
	return a.description
}

func (a *App) Meta() *metav1.ObjectMeta {
	return &a.ObjectMeta
}
