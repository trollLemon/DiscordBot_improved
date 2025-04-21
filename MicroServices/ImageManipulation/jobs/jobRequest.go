package jobs

import (
	"gocv.io/x/gocv"
)

type JobRequest struct {
	Job    *Job
	Result chan *gocv.Mat
	Error  chan error
}

func NewJobRequest(job *Job) *JobRequest {
	return &JobRequest{
		Job:    job,
		Result: make(chan *gocv.Mat, 1),
		Error:  make(chan error, 1),
	}
}
