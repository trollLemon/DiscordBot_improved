package jobs

import (
	"gocv.io/x/gocv"
)

type Operation interface {
	Run(input *gocv.Mat) (*gocv.Mat,  error)
}




/*
Job Struct

*/
type Job struct {
	jobId       uint8
	Operation
	inputImage  *gocv.Mat
}

func (j *Job) Process() (*gocv.Mat, error) {

	return j.Run(j.inputImage)
}


func (j *Job) GetJobId() uint8 {
	return j.jobId
}

type Invert struct{}

func (i Invert) Run(input gocv.Mat, output gocv.Mat) error {
	return nil
}


type EdgeDetect struct {
	t_lower  float64
	t_higher float64
}

func (e EdgeDetect) Run(input gocv.Mat, output gocv.Mat) error {
	return nil
}
