package service

import (
	"myapp/model"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"strconv"

	"k8s.io/client-go/kubernetes"
	batchv1 "k8s.io/api/batch/v1"
)
type KubernetesClient struct {
	clientset *kubernetes.Clientset
	queue ThreadSafePriorityQueue
	concurrency int
	ch chan string
}

func (client *KubernetesClient) Initialize(clientset *kubernetes.Clientset, concurrency int) {
	client.clientset = clientset
	client.queue = *NewThreadSafePriorityQueue()
	client.ch = make(chan string)
	client.concurrency = concurrency
	for range(concurrency) {
		go client.sender()
	}
	go client.receiver()
}

func (client *KubernetesClient) send_to_cluster(job batchv1.Job) bool{
	createdJob, err := client.clientset.BatchV1().Jobs(job.Namespace).Create(context.TODO(), &job, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Failed to create Job: %v", err)
		return false
	}

	fmt.Printf("Job created from file: %s/%s", createdJob.Namespace, createdJob.Name)
	return true
}

func (client *KubernetesClient) sender() {
	for {
		if client.queue.pq.Len() == 0 {
			continue
		} else {

			custom_job, ok := client.queue.Pop()

			if !ok {
				continue
			}
			res := client.send_to_cluster(custom_job.(*model.CustomJob).Job)
			client.ch <- strconv.FormatBool(res)
		}
	}
}

func (client *KubernetesClient) Submit(custom_job *model.CustomJob) {
	client.queue.Push(custom_job)
}

func (client *KubernetesClient) receiver() {
	for {
		response := <- client.ch
		fmt.Println(response)
	}
}