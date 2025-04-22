
package util

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/rest" // Import for in-cluster configuration
)

func Initialize_client() {
	// 1. Load Kubernetes configuration
	//    This is typically loaded from the kubeconfig file, which is the same file
	//    used by kubectl.  We'll try to load it from the default location,
	//    but you can also specify a specific path.
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		// Fallback: try the in-cluster configuration.  This is necessary
		// for when your application is running *inside* a Kubernetes
		// cluster.  It uses the service account credentials.
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		//sets global kubeconfig variable
		clientConfig := config
		fmt.Println(clientConfig)
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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

	// 4. Example: Creating a Pod (Optional)
	//    If you want to try creating a resource, you can uncomment this section.
	//    Make sure you have the necessary permissions in your cluster.
	/*
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-test-pod",
				Namespace: "default", // Change this if you want a different namespace
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx",
						Image: "nginx:1.23",
					},
				},
			},
		}

		createdPod, err := clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Created pod: %s in namespace: %s\n", createdPod.Name, createdPod.Namespace)
	*/

	fmt.Println("Done.")
}

