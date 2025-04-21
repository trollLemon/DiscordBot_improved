package worker

import (
	"github.com/rs/zerolog/log"
	"image_manip/jobs"
)

func Worker(workerId int, jobRequests <-chan *jobs.JobRequest) {

	for jobRequest := range jobRequests {

		job := jobRequest.Job

		log.Info().Msgf("Worker %d: Starting job: %d", workerId, job.GetJobId())

		result, err := job.Process()

		if err == nil {
			log.Info().Msgf("Worker %d: Finished job: %d in %d ns", workerId, job.GetJobId(), job.GetTimeElapsed())
			jobRequest.Result <- result
		} else {
			log.Error().Msgf("Worker %d: Error processing job %d: %v. Job started at %d and failed at %d", workerId, job.GetJobId(), err.Error(), job.GetStartTime().UnixMilli(), job.GetEndTime().UnixMilli())
			jobRequest.Error <- err
		}
	}
}
