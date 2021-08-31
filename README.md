[![Go Reference](https://pkg.go.dev/badge/github.com/kris-nova/naml.svg)](https://pkg.go.dev/github.com/kris-nova/naml)

# Not Another Markup Language.

> NAML is a Go library and command line tool that can be used as a framework to develop and deploy Kubernetes applications.

Replace Kubernetes YAML with raw Go!

Say so long 👋 to YAML and start using the Go 🎉 programming language to represent and deploy applications with Kubernetes.

Kubernetes applications are complicated, so lets use a proper Turing complete language to reason about them.

✅ Take advantage of all the lovely features of Go (Syntax highlighting, Cross compiling, Code generation, Documentation)

✅ Test your code directly in local Kubernetes using [kind](https://github.com/kubernetes-sigs/kind). Yes you can really deploy your applications to Kubernetes.

✅ Get your application directly into Go instead of YAML and use it in controllers, operators, CRs/CRDs easily. Use the Go compiler to your advantage.

## Convert YAML to Go

```bash
cat deploy.yaml | naml codify > main.go
```

Turn existing YAML into formatted and syntactically correct Go that implements the `Deployable` interface.

```bash
mkdir out

# Get started quickly with all objects in a namespace
kubectl get all -n default -o yaml | naml codify > out/main.go

# Pipe multiple .yaml files to a single Application
cat deployment.yaml service.yaml | naml codify \
  --author-name="Charlie" \
  --author-email="<charlie@nivenly.com>" > out/main.go
```

Copy the generic [Makefile](https://github.com/kris-nova/naml/blob/main/out/Makefile) to the same directory as your `main.go`

```bash 
wget https://raw.githubusercontent.com/kris-nova/naml/main/out/Makefile -o out/Makefile
cd out
make
./app -o yaml
./app install 
./app uninstall
```

Use `make help` for more. Happy coding 🎉.

## Example Projects

There is a "repository" of examples to borrow/fork:

- [values](https://github.com/naml-examples/full/blob/main/app.go#L52-L79) example pattern to replace a large `Values.yaml` file.
- [simple](https://github.com/naml-examples/simple) quick and simple example.
- [examples](https://github.com/naml-examples) GitHub organization.


### The Deployable Interface

As long as there is a Go system that implements this interface it can be used with `naml`. See examples for how to include an implementation in your project.

```go
// Deployable is an interface that can be implemented
// for deployable applications.
type Deployable interface {

    // Install will attempt to install in Kubernetes
    Install(client *kubernetes.Clientset) error

    // Uninstall will attempt to uninstall in Kubernetes
    Uninstall(client *kubernetes.Clientset) error

    // Meta returns the Kubernetes native ObjectMeta which is used to manage applications with naml.
    Meta() *metav1.ObjectMeta

    // Description returns the application description
    Description() string

    // Objects will return the runtime objects defined for each application
    Objects() []runtime.Object
}
```

In order to get the raw Kubernetes objects in Go without installing them anywhere, you pass in `nil` in place of an authenticated Kubernetes `Clientset`. 

Then you can access the objects in memory.

```go
    app.Install(nil)
    objects := app.Objects()
```

## Nothing fancy

There isn't anything special here. 🤷‍♀ We use the same client the rest of Kubernetes does.

 ❎ No new complex tools.

 ❎ No charts.

 ❎ No templating at runtime.

 ❎ No vague error messages.
 
 ❎ No more YAML guessing/checking.

 ✅ Just Go. 🎉

## Features

✨ There is not a single `.yaml` file in this entire repository. ✨

 - Express applications in 🎉 Go instead of YAML.
 - Use the Go compiler to check your syntax.
 - Write **real tests** 🤓 using Go to check and validate your deployments.
 - Test your applications in Kubernetes using [kind](https://github.com/kubernetes-sigs/kind).
 - Define custom installation logic. What happens if it fails? What about logical concerns at runtime?
 - Define custom application registries. Multiple apps of the same flavor? No problem.
 - Use the latest client (the same client the rest of Kubernetes uses).
