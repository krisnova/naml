package main

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/kris-nova/naml"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// main is the main entry point for your CLI application
func main() {

	naml.Register(&NAMLApp{"Example app 1"})
	naml.Register(&NAMLApp{"A great sample 2"})
	naml.Register(&NAMLApp{"Beeps Boops 3"})

	// Run the default CLI tooling
	err := naml.RunCommandLine()
	if err != nil {
		logger.Critical("%v", err)
		os.Exit(1)
	}
}

// NAMLApp is used for testing and debugging
type NAMLApp struct {
	Name string
}

func (n *NAMLApp) Install(client *kubernetes.Clientset) error {
	return nil
}

func (n *NAMLApp) Uninstall(client *kubernetes.Clientset) error {
	return nil
}

func (n *NAMLApp) Description() string {
	return "A wonderful sample application."
}

func (n *NAMLApp) Meta() *v1.ObjectMeta {
	return &v1.ObjectMeta{
		Name:            n.Name,
		ResourceVersion: "1.0.0",
	}
}
