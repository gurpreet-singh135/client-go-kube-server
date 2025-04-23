
package util

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest" // Import for in-cluster configuration
)

func Initialize_client(kubeconfig_path string) *kubernetes.Clientset{

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig_path)
	if err != nil {
        // If the default kubeconfig doesn't work, try in-cluster config.
        inClusterConfig, err := rest.InClusterConfig()
        if err != nil {
            fmt.Printf("Error loading kubeconfig: %v\n", err)
            fmt.Println("Falling back to in-cluster configuration failed as well.")
            os.Exit(1) // Exit if neither works.
        }
        config = inClusterConfig
        fmt.Println("Successfully used in-cluster configuration.")
	}

	// 2. Create the Kubernetes client
	//    This client is what we'll use to make requests to the Kubernetes API.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 3. Connect to the cluster and verify
	//    Let's do a simple test to make sure we can connect to the cluster.
	//    We'll list the namespaces in the cluster.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Successfully connected to Kubernetes cluster!")
	fmt.Println("Namespaces in the cluster:")
	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
	}

	fmt.Println("Done.")
	return clientset
}

