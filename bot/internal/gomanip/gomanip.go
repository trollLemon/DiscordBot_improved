package gomanip

import (

	"github.com/trollLemon/DiscordBot/internal/util"
)

func RandomFilter(gomanipClient *GoManip, image []byte, contentType string, kernelSize, lower, higher int64, normalize bool) ([]byte, error) {
	queries := util.RandomFilterQuery(kernelSize, lower, higher, normalize)
	bytes, err := gomanipClient.Do(image, contentType, "randomFilter", queries)
	return bytes, err 
}

func InvertImage(gomanipClient *GoManip, image []byte, contentType string) ([]byte, error) {
	bytes, err := gomanipClient.Do(image, contentType, "invert", "")
	return bytes, err
}

func SaturateImage(gomanipClient *GoManip, image []byte, contentType string, saturation int64) ([]byte, error) {
	saturationNorm := float32(saturation) / 100.0
	queries := util.SaturateQuery(saturationNorm)
	bytes, err := gomanipClient.Do(image, contentType, "saturate", queries)
	return bytes, err
}
func EdgeDetect(gomanipClient *GoManip, image []byte, contentType string, lower, higher int64) ([]byte, error) {
	queries := util.EdgeDetectQuery(lower, higher)
	bytes, err := gomanipClient.Do(image, contentType, "edgeDetection", queries)
	return bytes, err
}

func DilateImage(gomanipClient *GoManip, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	queries := util.DilateQuery(kernelSize, iterations)
	bytes, err := gomanipClient.Do(image, contentType, "morphology", queries)
	return bytes, err
}

func ErodeImage(gomanipClient *GoManip, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	queries := util.ErodeQuery(kernelSize, iterations)
	bytes, err := gomanipClient.Do(image, contentType, "morphology", queries)
	return bytes, err
}

func AddText(gomanipClient *GoManip, image []byte, contentType string, text string, fontScale float32, xPercentage float32, yPercentage float32) ([]byte, error) {
	queries := util.AddTextQuery(text, fontScale, xPercentage, yPercentage)
	bytes, err := gomanipClient.Do(image, contentType, "text", queries)
	return bytes, err
}

func Reduced(gomanipClient *GoManip, image []byte, contentType string, quality float32) ([]byte, error) {
	queries := util.ReduceQuery(quality)
	bytes, err := gomanipClient.Do(image, contentType, "reduced", queries)
	return bytes, err

}

func Shuffle(gomanipClient *GoManip, image []byte, contentType string, partitions int64) ([]byte, error) {
	queries := util.ShuffleQuery(partitions)
	bytes, err := gomanipClient.Do(image, contentType, "shuffle", queries)
	return bytes, err

}
