package worker

import (
	"context"
	"github.com/rs/zerolog/log"
	"goManip/jobs"
	"sync"
)

func Worker(shutdown context.Context, workerId int, jobRequests <-chan *jobs.JobRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	for jobRequest := range jobRequests {

		job := jobRequest.Job
		log.Info().Msgf("Worker %d: Starting job: %d", workerId, job.GetJobId())
		result, err := job.Process()
		select {

		case <-shutdown.Done():
			log.Info().Int("Worker", workerId).Msg("Worker shutdown")
			return

		// if the job timed out, no goroutine is waiting for a result.
		// we don't need to send anything
		case <-jobRequest.Ctx.Done():
			result.Close()

		default:

			if err != nil {
				log.Error().
					Str("Start time", job.GetStartTime().String()).
					Str("End time", job.GetEndTime().String()).
					Msgf("Worker %d Failed: %s", workerId, err.Error())
			} else {
				log.Info().
					Str("Start time", job.GetStartTime().String()).
					Str("End time", job.GetEndTime().String()).
					Int("Duration (ns)", job.GetTimeElapsed()).
					Msgf("Worker %d Completed", workerId)
			}

			jobRequest.Result <- &jobs.Result{Image: result, Error: err}

		}

	}
}
