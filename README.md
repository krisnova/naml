# YamYams

A Go framework for managing Kubernetes applications.

Say so long 👋 to YAML and start using the Go 🎉 programming language.

## About

This is a framework for infrastructure teams who need more than just conditional manifests.

Feel free to fork this repository and begin using it for your team. There isn't anything special here. 

## Adding a new application 

Copy the `_example` application and fill in the blanks. 

```bash 
cp -rv apps/_example apps/hello-world
```

Now edit the `apps/hello-world/app.go`.

```go
package helloworld

import (
	"fmt"
	yamyams "github.com/kris-nova/yamyams/pkg"
	"k8s.io/client-go/kubernetes"
)

type HelloWorld struct {
	meta      *yamyams.DeployableMeta
	resources []interface{}
}

func New() *HelloWorld {
	return &HelloWorld{
		meta: &yamyams.DeployableMeta{
			Name:        "Hello World Application!",
			Command:     "hello",
			Version:     "v0.0.1",
			Description: "This is a great hello world example!",
		},
	}
}

func (v *HelloWorld) Install(client *kubernetes.Clientset) error {
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hello-world",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "hello-world-container",
					Image: "busybox",
				},
			},
		},
	}
	newPod, err := client.CoreV1().Pods(v.namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	v.resources = append(v.resources, newPod)
	return nil
}

func (v *HelloWorld) Uninstall(client *kubernetes.Clientset) error {
	return client.CoreV1().Pods(v.namespace).Delete(context.TODO(), "hello-world", metav1.DeleteOptions{})
}

func (v *HelloWorld) Resources() []interface{} {
	return v.resources
}

func (v *HelloWorld) About() *yamyams.DeployableMeta {
	return v.meta
}
```

Now register your new app in `registry.go`

```go 
func Load() {
    // Register the new helloworld application
	yamyams.Register(helloworld.New())
}
```

Now compile the application to see if there are any errors. 

```bash 
go build -o yamyams cmd/*.go
```

Now you can install and uninstall in Kubernetes! 

```bash 
./yamyams install helloworld
./yamyams uninstall helloworld
```
