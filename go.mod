module github.com/kris-nova/naml

go 1.16

require (
	github.com/fatih/color v1.12.0
	github.com/hexops/valast v1.4.0
	github.com/kris-nova/logger v0.2.2
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	k8s.io/api v0.22.0
	k8s.io/apiextensions-apiserver v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
	sigs.k8s.io/kind v0.11.1
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/hexops/valast => github.com/fkautz/valast v1.4.1-0.20210806063143-f33a97256bcb
