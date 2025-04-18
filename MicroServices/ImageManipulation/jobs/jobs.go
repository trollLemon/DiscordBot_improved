package jobs

import (
	"time"

	"gocv.io/x/gocv"
)

type Operation interface {
	Run(input *gocv.Mat) (*gocv.Mat, error)
}

func NewJob(id uint8, operation Operation, image *gocv.Mat) *Job {

	return &Job{jobId: id, operation: operation, inputImage: image}
}

/*
Job Struct
*/
type Job struct {
	jobId uint8
	operation Operation
	inputImage *gocv.Mat
	startTime  time.Time
	endTime    time.Time
	elapsedTime time.Duration
}

func (j *Job) Process() (*gocv.Mat, error) {

	j.startTime = time.Now()
	result, err := j.operation.Run(j.inputImage)
	
	j.elapsedTime = time.Since(j.startTime)
	
	j.endTime = time.Now()

	return result, err
}

func (j *Job) GetJobId() uint8 {
	return j.jobId
}

func (j *Job) GetTimeElapsed() int {

	return int(j.elapsedTime)
}

func (j *Job) GetStartTime() int {
	return int(j.startTime.UnixMilli())
}

func (j *Job) GetEndTime() int {
	return int(j.endTime.UnixMilli())
}
