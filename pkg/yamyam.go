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

package yamyam

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/apps/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// YamYam is our very important business application
type YamYam struct {
	ContainerPort  int32
	ContainerImage string
	client         *kubernetes.Clientset
	deployment     *v1.Deployment
}

// New returns an empty *YamYam{}
func New() *YamYam {
	return &YamYam{}
}

// InstallKubernetes will try to install YamYam in Kubernetes
func (y *YamYam) InstallKubernetes() error {
	// Notice how we can do amazing things like add custom error messages
	// when things go wrong? This prevents us from digging through stack
	// traces later.
	//
	// The sky is the limit. We can make this method do anything we want
	// and return errors whenever something we don't like goes wrong.
	if y.client == nil {
		return fmt.Errorf("missing kube client: use KubernetesClient()")
	}
	// -----------------------------------------------------------------------------
	//
	// [beeps.yaml]
	//
	// If you can read YAML you can read this.
	// If you can edit YAML you can edit this.
	// If you can understand what {{ .Values.image }} does you can work on this.
	// If you can commit YAML changes you can commit Go changes.
	//
	// Oh did I mention that by doing it this way you never:
	//    - Have to guess if your indentation is wrong
	//    - Have to worry about typos (it wont compile)
	//    - Have to hunt around for errors (the go compiler tells you the line number)
	//    - Have to "render" a template to see what the fuck is going on
	//    - Have to wonder what happened at what time because you "changed" the template
	//
	// Furthermore if you use any of the readily available Go IDEs you get:
	//    - Tab hinting
	//    - Color coding
	//    - Syntax highlighting
	//    - Auto formatting
	//    - Auto linting
	//    - Documentation
	//    - Syntax checking
	//
	// Please stop writing YAML and just do this.
	//
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "yam-yams",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"beeps": "boops",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"beeps": "boops",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: "yam-yams",

							// ------------------------\
							Image: y.ContainerImage, // <-- THIS IS HOW YOU INTERPOLATE AT RUNTIME
							// ------------------------/

							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: y.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	// -----------------------------------------------------------------------------

	// Go to town.
	// Have fun.
	// Over engineer all kinds of crazy little systems and functions, methods, and types.
	// Build entire packages -- or -- Keep it simple.
	// Whatever you want to do.
	// We now have a turing complete programming language and the world is your oyster.
	//
	// See we can now do cool shit like validate our deployment actually looks the way we want
	//
	y.deployment = deployment
	err := y.Validate()
	if err != nil {
		// See how cool errors are? Especially when the errors are things that keep you from
		// getting paged.
		return fmt.Errorf("invalid deployment: %v", err)
	}
	// We could even make a copy of this and send it to MongoDB you know - just for record keeping.
	//
	// But why would you do that?
	//
	// Oh you know SO THAT WE HAVE A RECORD OF EVERYTHING OUR TEAM DOES SO WE CAN QUERY THE DATA LATER IF WE WANT
	//
	// The point is that you can build anything here because you are writing Go.
	//
	// You can introduce all kinds of cool features that make your team look awesome.
	err = y.Archive()
	if err != nil {
		return fmt.Errorf("unable to archive in mongodb: %s", err)
	}

	// And here we go folks
	//
	// The deployment is crafted.
	// The client is authenticated.
	// The validation checks are passed.
	// We can finally install in Kubernetes.
	result, err := y.client.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("oh no! something went wrong deploying to kubernetes: %v", err)
	}
	// Update the deployment in memory with the object from Kubernetes
	y.deployment = result
	return nil
}

func (y *YamYam) KubernetesClient(kubeconfigPath string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("unable to authenticate with Kubernetes with kube config %s: %v", kubeconfigPath, err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("unable to authenticate with Kubernetes with kube config %s: %v", kubeconfigPath, err)
	}
	y.client = client
	return nil
}

func (y *YamYam) Archive() error {
	// See you can even build out scaffolding for other features
	// you hope to build later. Maybe we will get to this next sprint.
	//
	// Because you are now a software engineer you do things like this.
	//
	// TODO @kris-nova save to mongo
	return nil
}

// Validate is all the handsome checks you get to design and talk about as a team.
// I wonder what would be important for you and your org to check for here?
// Hrmm...
func (y *YamYam) Validate() error {
	// Let's just make sure we have one...
	if y.deployment == nil {
		return fmt.Errorf("unable validate YamYam deployment: missing deployment")
	}
	// It wouldn't make much sense to try to deploy a deployment without any containers
	// Let's make sure we have at least 1 defined
	if len(y.deployment.Spec.Template.Spec.Containers) < 1 {
		return fmt.Errorf("unable to validate YamYam deployment. less than 1 container")
	}

	// So basically can do anything we want here because we have an entire programming language
	// at our disposal.
	//
	// Oh but you could use OPA and admissions controllers for things like this...
	//
	// Or we could just keep everything in one place and build a library of validation tools
	// like we would build a library of unit tests.
	//
	// You know... Fail quick.. Fail fast.. thats... a thing... right?
	//
	// Anyway I just dreamt up something we want to check for. Check that the name contains
	// the string "yam". Why? I don't know. Just seemed like a good example.
	//
	if !strings.Contains(y.deployment.Name, "yam") {
		return fmt.Errorf("unable validate YamYam deployment. invalid name %s", y.deployment.Name)
	}
	return nil
}

// UninstallKubernetes is a good method that is implemented poorly (on purpose).
// Feel free to come clean this up.
func (y *YamYam) UninstallKubernetes() error {
	if y.client == nil {
		return fmt.Errorf("missing kube client: use KubernetesClient()")
	}
	// This is actually really bad practice.
	// Here we "hard code" both the name of the namespace, and the name of the deployment to delete.
	// But you know what? A TODO is still 100% better than YAML files interpolated at runtime.
	// TODO @kris-nova introduce dynamic namespaces and names
	return y.client.AppsV1().Deployments("default").Delete(context.TODO(), "yam-yams", metav1.DeleteOptions{})

}

func int32Ptr(i int32) *int32 {
	return &i
}
