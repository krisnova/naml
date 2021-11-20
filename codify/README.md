# Codify

This is the `naml` feature that will convert YAML to Go.

### Adding an implementation 

Each implementation will require special attention to detail. 

The boilerplate is straightforward however there are some special considerations.

This is the `alias` method that will handle the weird Kubernetes import alias mechanism.

``` 
alias(generated, default string)
```

The main problem is that all objects will render as `v1` instead of their corresponding alias such as `corev1` or `metav1`.
This method will do it's best to figure these out (without using reflection).
There are members in the top of `codify.go` that call out the well-known subobjects that cannot be defaulted.
