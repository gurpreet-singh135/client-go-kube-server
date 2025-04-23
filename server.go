package main

import (
	"flag"
	"fmt"
	"myapp/handlers"
	"myapp/util"
	"os"
	"path/filepath"

)

func main() {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		home_dir = "/Users/gurpreet"
	}

	kubeconfig_path := flag.String("kubeconfig", filepath.Join(home_dir, ".kube", "config"), "the kubeconfig file path to connect to kubernetes cluster")
	server_port := flag.Int("port", 1323, "port on which application would run")
	max_concurrency := flag.Int("max-concurrency", 1, "maximum concurrency of submitting jobs to kubernetes cluster")


	flag.Parse()


	fmt.Println(*kubeconfig_path)
	fmt.Println(*max_concurrency)


	clientset := util.Initialize_client(*kubeconfig_path)

	if clientset == nil {
		panic("whoa")
	}


	e := handlers.Create_handlers(clientset, "default", *max_concurrency)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *server_port)))
}

