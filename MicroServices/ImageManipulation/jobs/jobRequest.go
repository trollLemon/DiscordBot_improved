package jobs

import (
	"context"
	"gocv.io/x/gocv"
)

type Result struct {
	Image *gocv.Mat
	Error error
}

type JobRequest struct {
	Job    *Job
	Result chan *Result
	Ctx    context.Context
}

func NewJobRequest(job *Job, ctx context.Context) *JobRequest {
	return &JobRequest{
		Job:    job,
		Result: make(chan *Result, 1),
		Ctx:    ctx,
	}
}
