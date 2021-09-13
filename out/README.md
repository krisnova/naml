# Out

This directory is for developing with NAML directly.

Use this directory to output generated Go code to. Then use the provided makefile to compile your `main.go` file. 

# Example 

Run this from the `naml` root directory.

```bash
kubectl run nginx --image nginx
kubectl get pods -o yaml | naml codify > out/main.go
cd out
make
./app 
```
