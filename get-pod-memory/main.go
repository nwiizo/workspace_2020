package main

import (
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {
	var kubeconfig string //empty, assuming inClusterConfig
	kubeconfig, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		panic(err)
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	podMetrics, err := mc.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(metav1.ListOptions{})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, podMetric := range podMetrics.Items {
		podContainers := podMetric.Containers
		for _, container := range podContainers {
			memQuantity, ok := container.Usage.Memory().AsInt64()
			if !ok {
				return
			}
			msg := fmt.Sprintf("Container Name: %s \n Memory usage: %d", container.Name, memQuantity)
			fmt.Println(msg)
		}
	}
}
