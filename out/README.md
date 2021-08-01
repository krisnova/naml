# Out

This directory is useful for local development of NAML.

```bash 
kubectl get deploy -oyaml | naml codify > out/main.go
cd out
make
make install
make clean
```