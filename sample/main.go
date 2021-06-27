package main

import (
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/logger"
	"github.com/kris-nova/naml"
	"k8s.io/client-go/kubernetes"
)

// main is the main entry point for your CLI application
func main() {

	naml.Register(&NAMLApp{"beeps"})
	naml.Register(&NAMLApp{"boops"})
	naml.Register(&NAMLApp{"meeps"})

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
	deployment := naml.BusyboxDeployment(n.Name)
	_, err := client.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func (n *NAMLApp) Uninstall(client *kubernetes.Clientset) error {
	return nil
}

func (n *NAMLApp) Description() string {
	return "A wonderful sample application."
}

func (n *NAMLApp) Meta() *metav1.ObjectMeta {
	return &metav1.ObjectMeta{
		Name:            n.Name,
		ResourceVersion: "1.0.0",
	}
}
