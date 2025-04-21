package handlers_test

import (
	"gocv.io/x/gocv"
	"image_manip/jobs"
	"image_manip/worker"
	"testing"
)

func jobsToRequests(jobArr []*jobs.Job) []*jobs.JobRequest {
	var jobRequests []*jobs.JobRequest

	for _, job := range jobArr {
		jobRequests = append(jobRequests, jobs.NewJobRequest(job))
	}

	return jobRequests
}

func TestJobDispatch(t *testing.T) {
	testImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV8UC3) // RGB image
	defer testImage.Close()
	numWorkers := 4
	requestsChan := make(chan *jobs.JobRequest, numWorkers)
	defer close(requestsChan)

	tests := []struct {
		name      string
		wantError bool
		jobReqs   []*jobs.JobRequest
	}{
		{
			name:      "TestJobDispatch Success",
			wantError: false,
			jobReqs: jobsToRequests([]*jobs.Job{
				jobs.NewJob(0, jobs.NewInvert(), &testImage),
				jobs.NewJob(1, jobs.NewSaturate(1.5), &testImage),
				jobs.NewJob(2, jobs.NewMorphology(3, 4, jobs.Dilate), &testImage),
				jobs.NewJob(3, jobs.NewMorphology(3, 10, jobs.Erode), &testImage),
				jobs.NewJob(4, jobs.NewReduce(0.2), &testImage),
				jobs.NewJob(5, jobs.NewEdgeDetection(120.0, 220.0), &testImage),
				jobs.NewJob(6, jobs.NewRandomFilter(9, -1, 1, false), &testImage),
				jobs.NewJob(7, jobs.NewShuffle(165), &testImage),
				jobs.NewJob(8, jobs.NewAddText("I love golang", 1.0, 0.5, 0.5), &testImage),
			},
			),
		},
		{
			name:      "TestJobDispatch Should Fail",
			wantError: true,
			jobReqs: jobsToRequests([]*jobs.Job{
				jobs.NewJob(1, jobs.NewSaturate(0.0), &testImage),
				jobs.NewJob(2, jobs.NewMorphology(-1, 4, jobs.Dilate), &testImage),
				jobs.NewJob(3, jobs.NewMorphology(3, 0, jobs.Erode), &testImage),
				jobs.NewJob(4, jobs.NewReduce(0.0), &testImage),
				jobs.NewJob(5, jobs.NewEdgeDetection(-10.0, 220.0), &testImage),
				jobs.NewJob(6, jobs.NewRandomFilter(0, -1, 1, false), &testImage),
				jobs.NewJob(7, jobs.NewShuffle(0), &testImage),
				jobs.NewJob(8, jobs.NewAddText("I love golang", 1.0, 2.0, 0.5), &testImage),
			},
			),
		},
	}

	for idx := range numWorkers {
		go worker.Worker(idx, requestsChan)
	}

	for _, tt := range tests {

		for _, jobReq := range tt.jobReqs {
			requestsChan <- jobReq
		}

		for _, jobReq := range tt.jobReqs {

			select {
			case result := <-jobReq.Result:
				if result == nil && !tt.wantError {
					t.Errorf("TestJobDispatch(): %s: jobReq result is nil, expected non-nil", tt.name)
				}
			case err := <-jobReq.Error:
				if tt.wantError && err == nil {
					t.Errorf("TestJobDispatch(): %s: got error: %v, expected error: %v", tt.name, err, tt.wantError)
				}
			}

		}

	}
}
