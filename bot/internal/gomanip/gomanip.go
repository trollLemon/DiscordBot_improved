package gomanip

import (
	"errors"
	"fmt"

	"github.com/trollLemon/DiscordBot/internal/apiErrors"
	"github.com/trollLemon/DiscordBot/internal/util"
)

var (
	ErrTimedOut  = errors.New("the service timed out")
	ErrBadParams = errors.New("the parameters are invalid")
	ErrGeneral   = errors.New("something went wrong, try the command again")
)

func errorChecker(err error) error {
	if errors.Is(err, apierrors.ErrRetry) {
		return ErrTimedOut
	}
	// this error contains the error response from the json, so we can have the error message here and let the user know.
	if errors.Is(err, apierrors.ErrAPI) {
		return fmt.Errorf("%v, %w", err, ErrBadParams)
	}
	if err != nil {
		return ErrGeneral 
	}
	return nil
}

func RandomFilter(gomanipClient *GoManip, image []byte, contentType string, kernelSize, lower, higher int64, normalize bool) ([]byte, error) {
	queries := util.RandomFilterQuery(kernelSize, lower, higher, normalize)
	bytes, err := gomanipClient.Do(image, contentType, "randomFilter", queries)
	return bytes, errorChecker(err)
}

func InvertImage(gomanipClient *GoManip, image []byte, contentType string) ([]byte, error) {
	bytes, err := gomanipClient.Do(image, contentType, "invert", "")
	return bytes, errorChecker(err)
}

func SaturateImage(gomanipClient *GoManip, image []byte, contentType string, saturation int64) ([]byte, error) {
	saturationNorm := float32(saturation) / 100.0
	queries := util.SaturateQuery(saturationNorm)
	bytes, err := gomanipClient.Do(image, contentType, "saturate", queries)
	return bytes, errorChecker(err)
}
func EdgeDetect(gomanipClient *GoManip, image []byte, contentType string, lower, higher int64) ([]byte, error) {
	queries := util.EdgeDetectQuery(lower, higher)
	bytes, err := gomanipClient.Do(image, contentType, "edgeDetection", queries)
	return bytes, errorChecker(err)
}

func DilateImage(gomanipClient *GoManip, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	queries := util.DilateQuery(kernelSize, iterations)
	bytes, err := gomanipClient.Do(image, contentType, "morphology", queries)
	return bytes, errorChecker(err)
}

func ErodeImage(gomanipClient *GoManip, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	queries := util.ErodeQuery(kernelSize, iterations)
	bytes, err := gomanipClient.Do(image, contentType, "morphology", queries)
	return bytes, errorChecker(err)
}

func AddText(gomanipClient *GoManip, image []byte, contentType string, text string, fontScale float32, xPercentage float32, yPercentage float32) ([]byte, error) {
	queries := util.AddTextQuery(text, fontScale, xPercentage, yPercentage)
	bytes, err := gomanipClient.Do(image, contentType, "text", queries)
	return bytes, errorChecker(err)
}

func Reduced(gomanipClient *GoManip, image []byte, contentType string, quality float32) ([]byte, error) {
	queries := util.ReduceQuery(quality)
	bytes, err := gomanipClient.Do(image, contentType, "reduced", queries)
	return bytes, errorChecker(err)

}

func Shuffle(gomanipClient *GoManip, image []byte, contentType string, partitions int64) ([]byte, error) {
	queries := util.ShuffleQuery(partitions)
	bytes, err := gomanipClient.Do(image, contentType, "shuffle", queries)
	return bytes, errorChecker(err)

}
