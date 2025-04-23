/*
   - /api/randomFilteredImage/{image_link:path}/{kernel_size}/{low}/{high}/{kernel_type}
   - /api/invertedImage/{image_link:path}
   - /api/saturatedImage/{image_link:path}/{saturation}
   - /api/edgeImage/{image_link:path}/{lower}/{higher}
   - /api/dilatedImage/{image_link:path}/{box_size}/{iterations}
   - /api/erodedImage/{image_link:path}/{box_size}/{iterations}
   - /api/textImage/{image_link:path}/{text}/{font_scale}/{x}/{y}
   - /api/reducedImage/{image_link:path}/{quality}
   - /api/shuffledImage/{image_link:path}/{partitions}
*/

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"image_manip/jobs"
	"image_manip/worker"
	"net/http"
)

func initRouting(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))

}

func GraceFullShutdown(jobReqs chan *jobs.JobRequest) {
	log.Info().Msg("Closing worker request channels")
	close(jobReqs)

}

func main() {
	e := echo.New()
	initRouting(e)
	numWorkers := 4
	jobReqs := make(chan *jobs.JobRequest, numWorkers)

	for workerId := range numWorkers {
		go worker.Worker(workerId, jobReqs)
	}

}
