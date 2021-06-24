# Simple Application

A basic example of how to build a naml project.

## app.go

Every project should define an `app.go` file.

The file should implement the `Deployable{}` interface.

```go
type Deployable interface {

    // Install will attempt to install in Kubernetes
    Install(client *kubernetes.Clientset) error

    // Uninstall will attempt to uninstall in Kubernetes
    Uninstall(client *kubernetes.Clientset) error

    // Meta returns the Kubernetes native ObjectMeta which is used to manage applications with naml.
    Meta() *v1.ObjectMeta
}
```

## app_test.go

We encourage everyone to also test their applications.

```go 
// TestAppName shows how you can test arbitrary parts of your application.
func TestAppName(t *testing.T) {
	app := New("default", "sample-app", "beeps-boops", 2)
	if app.Name != "sample-app" {
		t.Errorf(".Name is not plumbed through from New()")
	}
}
```

## cmd package

You can also take advantage of the default CLI tooling for `naml`.

```go
func main() {
	a := app.New("default", "simple-app", "beeps-boops", 17)
	naml.Register(a)
	err := cmd.RunCLI()
}
```

