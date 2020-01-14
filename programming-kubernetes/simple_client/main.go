package main

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var defaultKubeConfigPath string
	if home := homedir.HomeDir(); home != "" {
		// build kubeconfig path from $HOME dir
		defaultKubeConfigPath = filepath.Join(home, ".kube", "config")
	}

	// set kubeconfig flag
	kubeconfig := flag.String("kubeconfig", defaultKubeConfigPath, "kubeconfig config file")
	if kubeconfig == nil {
		fmt.Println("kubeconfig config error") // myError
		return
	}
	flag.Parse()

	// retrieve kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println("kubeconfig config error")
		return // myError
	}
	// get clientset for kubernetes resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("kubeconfig config error")
		return // myError
	}
	// Get list of pod objects
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("No pod")
		return // myError
	}
	// show pod object to stdout
	for i, pod := range pods.Items {
		fmt.Printf("[Pod Name %d]%s\n", i, pod.GetName())
	}
}
