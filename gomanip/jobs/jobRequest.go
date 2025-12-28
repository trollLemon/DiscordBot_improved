package jobs

import (
	"context"
	"gocv.io/x/gocv"
)

// Result stores the Result image (which may be nil if there is an error), and an error if one occured. 
type Result struct {
	Image *gocv.Mat
	Error error
}

// JobRequest facilitates creating a job and listening on a channel for completion. 
type JobRequest struct {
	Job    *Job
	Result chan *Result
	Ctx    context.Context
}

// NewJobRequest returns a new JobRequest.
func NewJobRequest(job *Job, ctx context.Context) *JobRequest {
	return &JobRequest{
		Job:    job,
		Result: make(chan *Result, 1),
		Ctx:    ctx,
	}
}
