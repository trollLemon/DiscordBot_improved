package server

import (
	"context"

	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"

	"goManip/jobdispatch"
	"goManip/jobs"
	"goManip/util"
)





func getDispatcher(c echo.Context) *jobdispatch.JobDispatcher {
	return c.Get("jobDispatcher").(*jobdispatch.JobDispatcher)
}

func handleImageOperation(
	c echo.Context,
	processFunc func(image *gocv.Mat) (*gocv.NativeByteBuffer, error),
) error {

	image, err := util.GetImageFromBody(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read image")
		return SendGomanipError(c, err)
	}

	resultImage, err := processFunc(image)
	if err != nil {
		log.Error().Err(err).Msg("Image processing failed")
		return SendGomanipError(c, err)
	}
	defer resultImage.Close()

	return c.Blob(http.StatusOK, "image/png", resultImage.GetBytes())
}

func InvertEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueInvertImage(jobDispatcher, image)
	})

}

func SaturateEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	saturation, err := util.ParseSaturation(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse saturation")
		return c.String(http.StatusBadRequest, "Failed to parse saturation: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueSaturateImage(jobDispatcher, image, saturation)
	})

}

func EdgeDetectionEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	tLower, tHigher, err := util.ParseEdgeDetection(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse edge detection")
		return c.String(http.StatusBadRequest, "Failed to parse edge detection: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueDetectEdges(jobDispatcher, image, tLower, tHigher)
	})
}

func MorphologyEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	morphType, kernelSize, iterations, err := util.ParseMorphology(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse morphology")
		return c.String(http.StatusBadRequest, "Failed to parse morphology: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Choice(morphType), kernelSize, iterations)
	})
}

func ReduceEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	quality, err := util.ParseReduce(c)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse reduce")
		return c.String(http.StatusBadRequest, "Failed to parse reduce: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueReduceImage(jobDispatcher, image, quality)
	})
}

func AddTextEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	text, fontScale, xPerc, yPerc, err := util.ParseAddText(c)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse add text")
		return c.String(http.StatusBadRequest, "Failed to parse add text: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueAddText(jobDispatcher, image, text, fontScale, xPerc, yPerc)

	})
}

func RandomFilterEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	minVal, maxVal, kernelSize, normalize, err := util.ParseRandomFilter(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse random filter")
		return c.String(http.StatusBadRequest, "Failed to parse random filter: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueRandomFilter(jobDispatcher, image, minVal, maxVal, kernelSize, normalize)
	})
}
func ShuffleEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		log.Error().Msg("Job dispatcher is not present in the context")
		return c.String(http.StatusInternalServerError, "failed to get job dispatcher")
	}

	partitions, err := util.ParseShuffle(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse shuffle")
		return c.String(http.StatusBadRequest, "Failed to parse shuffle: "+err.Error())
	}

	return handleImageOperation(c, func(image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
		return jobdispatch.EnqueueShuffle(jobDispatcher, image, partitions)
	})
}

func InitRouting(e *echo.Echo, jobDispatcher *jobdispatch.JobDispatcher, address, port string) {
	e.Use(JobDispatcherMiddleware(jobDispatcher))
	e.Use(FileTypeVerifyMiddleware())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogURIPath: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Str("Path", v.URIPath).
				Int("status", v.Status).
				Str("method", c.Request().Method).
				Str("Time", v.StartTime.String()).
				Msg("Request")
			return nil
		},
	}))

	e.POST("/invert/", InvertEndpoint)
	e.POST("/saturate/", SaturateEndpoint)
	e.POST("/edgeDetection/", EdgeDetectionEndpoint)
	e.POST("/morphology/", MorphologyEndpoint)
	e.POST("/reduction/", ReduceEndpoint)
	e.POST("/text/", AddTextEndpoint)
	e.POST("/randomFilter/", RandomFilterEndpoint)
	e.POST("/shuffle/", ShuffleEndpoint)
	e.Logger.Fatal(e.Start(address + ":" + port))

}

func GraceFullShutdown(jobDispatcher *jobdispatch.JobDispatcher, wg *sync.WaitGroup, cancel context.CancelFunc) {
	log.Info().Msg("Closing worker request channels")
	jobDispatcher.Close()
	log.Info().Msg("Stopping Workers")
	cancel()
	wg.Wait()

}
