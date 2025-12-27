package jobs

import (
	"gocv.io/x/gocv"
	"time"
)

// Operation represents an abstraction for some operation to perform on an image.
// A struct implementing Operation should have any parameters stored within the struct, so the Job struct can call Run and get a 
// result image. 
type Operation interface {
	Run(input *gocv.Mat) (*gocv.Mat, error)
}

// NewJob creates a new job given an id, operation, and input image.
func NewJob(id uint32, operation Operation, image *gocv.Mat) *Job {

	return &Job{jobId: id, operation: operation, inputImage: image}
}

// Job contains the input image, operation object, and a jobId, as well as startTime and endTime which are populated after a worker calls Process().
type Job struct {
	jobId       uint32
	operation   Operation
	inputImage  *gocv.Mat
	startTime   time.Time
	endTime     time.Time
	elapsedTime time.Duration
}

// Process runs the operation on the stored input image, and returns the result image and error (if any)
func (j *Job) Process() (*gocv.Mat, error) {
	j.startTime = time.Now()

	result, err := j.operation.Run(j.inputImage)

	j.elapsedTime = time.Since(j.startTime)

	j.endTime = time.Now()

	return result, err
}

// GetJobId returns the job id.
func (j *Job) GetJobId() uint32 {
	return j.jobId
}

// GetTimeElapsed returns the total amount of time (in ns) the job took.
func (j *Job) GetTimeElapsed() int {
	return int(j.elapsedTime)
}

// GetStartTime returns the time the job started.
func (j *Job) GetStartTime() time.Time {
	return j.startTime
}

// GetEndTime returns the time the job ended.
func (j *Job) GetEndTime() time.Time {
	return j.endTime
}
