# Applications

Define your applications here.

Every application should implement the `Deployable` interface.

An application can also be defined in any repository as long as the `Deployable` interface is implemented it can be used with the rest of naml.

### Using Go instead of YAML templating for a Pod

Traditionally we would template YAML doing something like this.

```yaml
  apiVersion: v1
  kind: Pod
    labels:
      app: beeps-boops
    name: {{ .Values.Name }}
    namespace: default
  spec:
    containers:
      image: {{ .Values.Image }}
      name: {{ .Values.Name }}
```

This presented many problems, some of the most obvious being that the code was loosely typed, a lack of validation and testing, and syntax checking.

The Go programming language and most modern IDEs gives us a lot of features we can use to our advantage.

In this example we "plumb" values (and their types) through to our code.

```go 
type MyPod struct {
	resources     []interface{}
	meta          *naml.DeployableMeta
	namespace     string
	name          string
	image         string
	exampleInt    int
}

func New(namespace string, name string, image string, exampleInt int) *MyPod {
	return &MyPod{
		namespace:     namespace,
		exampleInt:    exampleInt,
		image:         image,
		name:          name,
		meta: &naml.DeployableMeta{
			Name:        "Example Pod",
			Version:     "0.0.1",
			Command:     "mypod",
			Description: "A simple example Kubernetes pod",
		},
	}
}
```

Now we can implement logic specifically to replace the features we had with templating YAML, and harden our applications even more.

```go 
func (v *MyPod) Install(client *kubernetes.Clientset) error {

	// Check if the name contains a substring
	if !strings.Contains(v.name, "example") {
		return fmt.Errorf("invalid name %s, must contain substring 'example'", v.name)
	}

	// Check if exampleInt is greater than 1
	if v.exampleInt < 2 {
		return fmt.Errorf("invalid exampleInt %d, must be greater than 1", exampleInt)
	}

	// Now we can "plumb" values through to our pod
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: v.name,
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  v.name,
					Image: v.image,
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
```


