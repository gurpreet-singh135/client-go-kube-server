package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"

	"context"
	"fmt"
	"io"
	"myapp/model"
	"myapp/service"

	"strconv"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func Create_handlers(clientset *kubernetes.Clientset, namespace string, concurrency int) *echo.Echo{

	var kubernetes_client *service.KubernetesClient = &service.KubernetesClient{};
	kubernetes_client.Initialize(clientset, concurrency)
	e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	e.POST("/jobs", func(c echo.Context) error {
		return createJobFromFileHandler(c, kubernetes_client)
	})
	
	e.GET("/jobs/pending", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! /jobs")
	})

	e.GET("/jobs/running", func(c echo.Context) error {
		jobList, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println("Failed to list Jobs: %v", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to list Jobs: %v", err))
		}
		jobs := make([]string, 0)
		for _, job := range jobList.Items {
			if job.Status.Active > 0 {
				jobs = append(jobs, job.Name)
			}
		}

		return c.String(http.StatusOK, fmt.Sprintf("running jobs are: %v", jobs))
	})

	return e
	
}

func createJobFromFileHandler(c echo.Context, kubernetes_client *service.KubernetesClient) error {

	priority_s := c.FormValue("priority") // "jobFile" is the name of the form field
	if priority_s == ""{
		fmt.Printf("Failed to get priority")
		priority_s = "1000"
	}
	priority, err := strconv.Atoi(priority_s)

	if err != nil {
		fmt.Printf("Failed to get concurrency: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing concurrency"})
    }
	// 1. Get the uploaded file.
	file, err := c.FormFile("jobFile") // "jobFile" is the name of the form field
	if err != nil {
		fmt.Printf("Failed to get uploaded file: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'jobFile' form field or file"})
	}
	src, err := file.Open()
	if err != nil {
		fmt.Printf("Failed to open uploaded file: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	// 2. Read the file content.  We limit the size to prevent malicious uploads.
	maxBytes := int64(2 * 1024 * 1024) // 2MB max
	reader := io.LimitReader(src, maxBytes)
	fileBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read uploaded file"})
	}

	// 3. Unmarshal the YAML into a Job object.  Use the Kubernetes YAML library.
	var job batchv1.Job
	if err := yaml.Unmarshal(fileBytes, &job); err != nil {
		fmt.Printf("Failed to unmarshal YAML: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Job YAML: " + err.Error()})
	}
    if job.Namespace == "" {
        job.Namespace = "default" //set to default if namespace not provided in request.
    }

	kubernetes_client.Submit(&model.CustomJob{
		Job: job,
		Priority: priority,
		Index: 0,
	})



	fmt.Printf("Job created from file: %s/%s", job.Namespace, job.Name)
	return c.JSON(http.StatusCreated, job)
}
