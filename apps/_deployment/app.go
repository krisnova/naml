//
// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, softwar
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//    ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗
//    ████╗  ██║██╔═████╗██║   ██║██╔══██╗
//    ██╔██╗ ██║██║██╔██║██║   ██║███████║
//    ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║
//    ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║
//    ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝

package mydeployment

import (
	"context"
	"fmt"
	yamyams "github.com/kris-nova/yamyams/pkg"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"

	v1 "k8s.io/api/apps/v1"

	apiv1 "k8s.io/api/core/v1"
)

// MyDeployment represents a Kubernetes deployment.
type MyDeployment struct {
	resources     []interface{}
	meta          *yamyams.DeployableMeta
	namespace     string
	name          string
	exampleString string
	exampleInt    int
}

// New will return a new MyDeployment.
//
// We can pass in custom arguments here and use them when we Install().
// In this example we use exmpleString as a label, and exampleInt as our replica count.
func New(namespace string, name string, exampleString string, exampleInt int) *MyDeployment {
	return &MyDeployment{
		namespace:     namespace,
		exampleInt:    exampleInt,
		exampleString: exampleString,
		name:          name,
		meta: &yamyams.DeployableMeta{
			Name:        "Example Deployment",
			Version:     "0.0.1",
			Command:     "mydeployment",
			Description: "A simple example Kubernetes deployment",
		},
	}
}

// Install will try to install in Kubernetes
func (v *MyDeployment) Install(client *kubernetes.Clientset) error {

	labels := map[string]string{
		"k8s-app":       "mydeployment",
		"app":           "mydeployment",
		"example-label": v.exampleString,
	}

	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: v.name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: yamyams.I32p(int32(v.exampleInt)),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  v.name,
							Image: "busybox",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	updatedDeployment, err := client.AppsV1().Deployments(v.namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to install deployment in Kubernetes: %v", err)
	}
	v.resources = append(v.resources, updatedDeployment)
	return nil
}

func (v *MyDeployment) Uninstall(client *kubernetes.Clientset) error {
	return client.AppsV1().Deployments(v.namespace).Delete(context.TODO(), v.name, metav1.DeleteOptions{})
}

func (v *MyDeployment) Resources() []interface{} {
	return v.resources
}

func (v *MyDeployment) About() *yamyams.DeployableMeta {
	return v.meta
}
