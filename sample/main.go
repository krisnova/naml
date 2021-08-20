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
	"os"

	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/hexops/valast"
	"github.com/kris-nova/naml"
	"k8s.io/client-go/kubernetes"
)

var Version string = "1.2.3"

func main() {
	// Load the application into the NAML registery
	// Note: naml.Register() can be used multiple times.
	//
	naml.Register(NewApp("App-Barnaby", "App for Barnaby!"))
	naml.Register(NewApp("App-Nova", "App for Nova!"))

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
	objects     []runtime.Object
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
			Name: name,
		},
		// --------------------
		// Add your fields here
		// --------------------
	}
}

func (a *App) Install(client *kubernetes.Clientset) error {
	var err error

	// Adding a deployment: "nginx"
	nginxDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/appsv1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "nginx",
			Namespace:   "default",
			UID:         types.UID("c39ccefc-491a-4857-bb47-1aa540f2129a"),
			Generation:  1,
			Labels:      map[string]string{"app": "nginx"},
			Annotations: map[string]string{"deployment.kubernetes.io/revision": "1"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: valast.Addr(int32(1)).(*int32),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"app": "nginx",
			}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "nginx"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{corev1.Container{
						Name:                     "nginx",
						Image:                    "nginx",
						TerminationMessagePath:   "/dev/termination-log",
						TerminationMessagePolicy: corev1.TerminationMessagePolicy("File"),
						ImagePullPolicy:          corev1.PullPolicy("Always"),
					}},
					RestartPolicy:                 corev1.RestartPolicy("Always"),
					TerminationGracePeriodSeconds: valast.Addr(int64(30)).(*int64),
					DNSPolicy:                     corev1.DNSPolicy("ClusterFirst"),
					SecurityContext:               &corev1.PodSecurityContext{},
					SchedulerName:                 "default-scheduler",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType("RollingUpdate"),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Type(1),
						StrVal: "25%",
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Type(1),
						StrVal: "25%",
					},
				},
			},
		},
	}

	if client != nil {
		_, err = client.AppsV1().Deployments("default").Create(context.TODO(), nginxDeployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// Adding a deployment: "nginx-deployment"
	nginx_deploymentDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/appsv1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "nginx-deployment",
			Namespace:  "default",
			UID:        types.UID("445061d9-5000-471b-8e06-45f5240dedb6"),
			Generation: 3,
			Annotations: map[string]string{
				"deployment.kubernetes.io/revision": "3",
				"kubectl.kubernetes.io/last-applied-configuration": `{"apiVersion":"apps/appsv1","kind":"Deployment","metadata":{"annotations":{},"name":"nginx-deployment","namespace":"default"},"spec":{"replicas":4,"selector":{"matchLabels":{"app":"nginx"}},"template":{"metadata":{"labels":{"app":"nginx"}},"spec":{"containers":[{"image":"nginx:1.14.2","name":"nginx","ports":[{"containerPort":80}]}]}}}}
`,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: valast.Addr(int32(4)).(*int32),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "nginx"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{
					"app": "nginx",
				}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{corev1.Container{
						Name:  "nginx",
						Image: "nginx:1.14.2",
						Ports: []corev1.ContainerPort{
							corev1.ContainerPort{
								ContainerPort: 80,
								Protocol:      corev1.Protocol("TCP"),
							},
						},
						TerminationMessagePath:   "/dev/termination-log",
						TerminationMessagePolicy: corev1.TerminationMessagePolicy("File"),
						ImagePullPolicy:          corev1.PullPolicy("IfNotPresent"),
					}},
					RestartPolicy:                 corev1.RestartPolicy("Always"),
					TerminationGracePeriodSeconds: valast.Addr(int64(30)).(*int64),
					DNSPolicy:                     corev1.DNSPolicy("ClusterFirst"),
					SecurityContext:               &corev1.PodSecurityContext{},
					SchedulerName:                 "default-scheduler",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType("RollingUpdate"),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Type(1),
						StrVal: "25%",
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Type(1),
						StrVal: "25%",
					},
				},
			},
		},
	}

	a.objects = append(a.objects, nginx_deploymentDeployment)
	if client != nil {
		_, err = client.AppsV1().Deployments("default").Create(context.TODO(), nginx_deploymentDeployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return err
}

func (a *App) Uninstall(client *kubernetes.Clientset) error {
	var err error

	err = client.AppsV1().Deployments("default").Delete(context.TODO(), "nginx", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = client.AppsV1().Deployments("default").Delete(context.TODO(), "nginx-deployment", metav1.DeleteOptions{})
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

func (a *App) Objects() []runtime.Object {
	return a.objects
}
