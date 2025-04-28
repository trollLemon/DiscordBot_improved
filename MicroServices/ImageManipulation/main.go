package main

import (
	"context"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"goManip/JobDispatch"
	"goManip/jobs"
	"goManip/util"
	"goManip/worker"
	"gocv.io/x/gocv"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func getDispatcher(c echo.Context) *JobDispatch.JobDispatcher {
	return c.Get("jobDispatcher").(*JobDispatch.JobDispatcher)
}

func handleImageOperation(
	c echo.Context,
	processFunc func(image *gocv.Mat) (*gocv.NativeByteBuffer, error),
) error {

	image, err := util.GetImageFromBody(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read image")
		return c.String(http.StatusBadRequest, "Failed to read image: "+err.Error())
	}

	resultImage, err := processFunc(image)
	if err != nil {
		log.Error().Err(err).Msg("Image processing failed")
		return c.String(http.StatusBadRequest, "Image processing failed: "+err.Error())
	}
	defer resultImage.Close()

	return c.Blob(http.StatusOK, "image/png", resultImage.GetBytes())
}

func InvertEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return JobDispatch.EnqueueInvertImage(jobDispatcher, image)
	})

}

func SaturateEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	saturation, err := util.ParseSaturation(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse saturation")
		return c.String(http.StatusBadRequest, "Failed to parse saturation: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return JobDispatch.EnqueueSaturateImage(jobDispatcher, image, saturation)
	})

}

func EdgeDetectionEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	tLower, tHigher, err := util.ParseEdgeDetection(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse edge detection")
		return c.String(http.StatusBadRequest, "Failed to parse edge detection: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return JobDispatch.EnqueueDetectEdges(jobDispatcher, image, tLower, tHigher)
	})
}

func MorphologyEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	morphType, kernelSize, iterations, err := util.ParseMorphology(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse morphology")
		return c.String(http.StatusBadRequest, "Failed to parse morphology: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return JobDispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Choice(morphType), kernelSize, iterations)
	})
}

func ReduceEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	quality, err := util.ParseReduce(c)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse reduce")
		return c.String(http.StatusBadRequest, "Failed to parse reduce: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return JobDispatch.EnqueueReduceImage(jobDispatcher, image, quality)
	})
}

func AddTextEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	text, fontScale, xPerc, yPerc, err := util.ParseAddText(c)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse add text")
		return c.String(http.StatusBadRequest, "Failed to parse add text: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {

		return JobDispatch.EnqueueAddText(jobDispatcher, image, text, fontScale, xPerc, yPerc)

	})
}

func RandomFilterEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	minVal, maxVal, kernelSize, normalize, err := util.ParseRandomFilter(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse random filter")
		return c.String(http.StatusBadRequest, "Failed to parse random filter: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return JobDispatch.EnqueueRandomFilter(jobDispatcher, image, minVal, maxVal, kernelSize, normalize)
	})
}
func JobDispatcherMiddleware(jobDispatcher *JobDispatch.JobDispatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("jobDispatcher", jobDispatcher)
			return next(c)
		}
	}
}

func ShuffleEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)

	partitions, err := util.ParseShuffle(c)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse shuffle")
		return c.String(http.StatusBadRequest, "Failed to parse shuffle: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return JobDispatch.EnqueueShuffle(jobDispatcher, image, partitions)
	})
}

func initRouting(e *echo.Echo, jobDispatcher *JobDispatch.JobDispatcher) {

	e.Use(JobDispatcherMiddleware(jobDispatcher))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("method", c.Request().Method).
				Str("Time", v.StartTime.String()).
				Msg("Request")
			return nil
		},
	}))

	e.POST("/api/image/invert/", InvertEndpoint)
	e.POST("/api/image/saturate/", SaturateEndpoint)
	e.POST("/api/image/edgeDetection/", EdgeDetectionEndpoint)
	e.POST("/api/image/morphology/", MorphologyEndpoint)
	e.POST("/api/image/reduction/", ReduceEndpoint)
	e.POST("/api/image/text/", AddTextEndpoint)
	e.POST("/api/image/randomFilter/", RandomFilterEndpoint)
	e.POST("/api/image/shuffle/", ShuffleEndpoint)
	e.Logger.Fatal(e.Start(":8080"))

}

func GraceFullShutdown(jobDispatcher *JobDispatch.JobDispatcher, wg *sync.WaitGroup, cancel context.CancelFunc) {

	log.Info().Msg("Closing worker request channels")
	jobDispatcher.Close()
	log.Info().Msg("Stopping Workers")
	cancel()
	wg.Wait()

}

func main() {
	prettyPrint := flag.Bool("pretty_print", false, "Enable pretty print output")
	numWorkers := flag.Int("num_workers", runtime.NumCPU(), "Number of workers")
	flag.Parse()

	if *prettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	jobReqs := make(chan *jobs.JobRequest, *numWorkers)
	maxTime := time.Second * 10
	jobDispatcher := JobDispatch.NewJobDispatcher(jobReqs, maxTime)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		GraceFullShutdown(jobDispatcher, wg, cancel)
		os.Exit(0)
	}()

	for workerId := range *numWorkers {
		log.Info().Msgf("Starting worker #%d", workerId+1)
		wg.Add(1)
		go worker.Worker(ctx, workerId+1, jobReqs, wg)
	}
	e := echo.New()
	initRouting(e, jobDispatcher)
}
