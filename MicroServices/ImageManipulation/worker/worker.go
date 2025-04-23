package worker

import (
	"context"
	"github.com/rs/zerolog/log"
	"image_manip/jobs"
	"sync"
)

func Worker(ctx context.Context, workerId int, jobRequests <-chan *jobs.JobRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	for jobRequest := range jobRequests {

		job := jobRequest.Job

		log.Info().Msgf("Worker %d: Starting job: %d", workerId, job.GetJobId())

		result, err := job.Process()
		select {

		case <-ctx.Done():
			log.Info().Msgf("Worker %d: shutdown", workerId)
			return

		default:
			if err == nil {
				log.Info().Msgf("Worker %d: Finished job: %d in %d ns", workerId, job.GetJobId(), job.GetTimeElapsed())
				jobRequest.Result <- result
				close(jobRequest.Result)
			} else {
				log.Error().Msgf("Worker %d: Error processing job %d: %v. Job started at %d and failed at %d", workerId, job.GetJobId(), err.Error(), job.GetStartTime().UnixMilli(), job.GetEndTime().UnixMilli())
				jobRequest.Error <- err
				close(jobRequest.Error)
			}
		}
	}
}
