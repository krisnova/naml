# YamYams

Replace Kubernetes YAML with raw Go!

Say so long 👋 to YAML and start using the Go 🎉 programming language.

Take advantage of all the lovely features of Go.

Test your code directly in local Kubernetes using [kind](https://github.com/kubernetes-sigs/kind).

## About

This is a framework for infrastructure teams who need more than just conditional manifests. 

This allows teams to start encapsulating, managing, and testing their applications in raw Go.

Teams can now buid controllers, operators, and custom toolchains using reliable, testable, and scalable Go.

## Nothing fancy

Feel free to fork this repository and begin using it for your team. There isn't anything special here. 🤷‍♀ We use the same client the rest of Kubernetes does.

 ❎ No new tools.

 ❎ No charts.

 ❎ No barrier to entry from a new tool.

 ❎ No templating at runtime.

 ❎ No vague error messages and counting YAML indentation.

 ✅ Just plain Go.

## Features

✨ There is not a single `.yaml` file in this entire repository. ✨

 - Express applications in 🎉 Go instead of YAML.
 - Use the Go compiler to check your syntax.
 - Write **real tests** 🤓 using Go to check and validate your deployments.
 - Test your applications in Kubernetes using [kind](https://github.com/kubernetes-sigs/kind).
 - Define custom installation logic. What happens if it fails?
 - Define custom application registries. Multiple apps of the same flavor? No problem.
 - Use the latest client (the same client the rest of Kubernetes uses).

## Creating a new app 

Copy `sampleapp` to a new folder in `/apps` to get started.

```bash 
cp -rv apps/sampleapp apps/hello
```

Edit your new application `/apps/hello/app.go`.

Use the same [client](https://github.com/kubernetes/client-go) the rest of Kubernetes uses to express your application.

Register your application in `/registry.go`.


```go
func Load(){
    yamyams.Register(hello.New("default", "hello-app"))
}
```



## Testing your application 

Create a new test for `hello` and edit the tests.

```bash 
cp apps/sampleapp_test.go apps/hello_test.go
```

Then bootstrap a local kind cluster and actually deploy your application to a real Kubernetes cluster.

```bash 
sudo -E make test
```

You can add custom integration tests to check whatever you want.

```go 
func TestHelloAppName(t *testing.T) {
	app := hello.New("default", "hello-app")
	if app.Name != "hello-app" {
		t.Errorf("app Name is not plumbed through from New()")
	}
}
```

## Deploy your application 

You can now deploy your application directly to Kubernetes from the provided CLI tool.

```bash 
make
sudo -E make install
yamyams install hello-app
```

You can also `list` and `uninstall`

```bash 
yamyams list
yamyams uninstall
```

## Use your application in many ways.

Now that your application is expressed in Go you can easily use tools like Kubernetes controllers and CRDs to manage and reconcile your application.

 - [KubeBuilder](https://github.com/kubernetes-sigs/kubebuilder) can help you build CRDs and operators.
 - [Operator Framework](https://github.com/operator-framework/operator-sdk) can help you build CRDs and operators.
 - [cdk8s](https://github.com/cdk8s-team/cdk8s) can be used to generate YAML in a similar way this project represents applications.