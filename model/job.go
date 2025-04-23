package model

import(
	batchv1 "k8s.io/api/batch/v1"
)

type CustomJob struct {
	Job batchv1.Job
	Priority int
	Index int
}