[![Go Reference](https://pkg.go.dev/badge/github.com/kris-nova/naml.svg)](https://pkg.go.dev/github.com/kris-nova/naml)

# Not another markup language. 

Replace Kubernetes YAML with raw Go!

Say so long 👋 to YAML and start using the Go 🎉 programming language to represent and deploy applications.

Take advantage of all the lovely features of Go.

Test your code directly in local Kubernetes using [kind](https://github.com/kubernetes-sigs/kind).

Get your application directly into Go instead of YAML and use it in controllers, operators, CRs/CRDs easily. Use the Go compiler to your advantage.

### Okay, but seriously

This is just `client-go`. There isn't really anything more to see here.

#### Quickstart

Check out the [examples](https://github.com/naml-examples) GitHub organization.

- [simple](https://github.com/naml-examples/simple) basic CLI example.
- [full](https://github.com/naml-examples/full) CLI example with custom commands and flags.

#### Implement Deployable

As long as there is a Go system that implements this interface it can be used with `naml`. See examples for how to include an implementation in your project.

```go
// Deployable is used to deploy applications.
// v0.2.0
type Deployable interface {

	// Install will attempt to install in Kubernetes
	Install(client *kubernetes.Clientset) error

	// Uninstall will attempt to uninstall in Kubernetes
	Uninstall(client *kubernetes.Clientset) error

	// Meta returns the Kubernetes native ObjectMeta which is used to manage applications with naml.
	Meta() *v1.ObjectMeta
	
	// Description will return the description of the application.
	Description() string
}
```

## About

This is a framework for infrastructure teams who need more than just conditional manifests. 

This allows teams to start encapsulating, managing, and testing their applications in raw Go.

Teams can now buid controllers, operators, and custom toolchains using reliable, testable, and scalable Go.

## The philosophy

The bet here is that any person confident in managing `YAML` for Kubernetes can also be equally as confident managing Go for Kubernetes.

The feature add is that no matter how good our YAML management tools get, they will never be as good as just plain Go when it comes to things like syntax checking, testing, shipping, and flexibility. 

## Nothing fancy

Feel free to fork this repository and begin using it for your team. There isn't anything special here. 🤷‍♀ We use the same client the rest of Kubernetes does.

 ❎ No new tools.

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
 - Define custom installation logic. What happens if it fails?
 - Define custom application registries. Multiple apps of the same flavor? No problem.
 - Use the latest client (the same client the rest of Kubernetes uses).


## Getting Started

Check out the [examples](https://github.com/naml-examples) GitHub organization. 

- [simple](https://github.com/naml-examples/simple) basic CLI example.
- [full](https://github.com/naml-examples/full) CLI example with custom commands and flags.
