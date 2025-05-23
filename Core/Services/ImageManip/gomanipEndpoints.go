package imagemanip

import "bot/util"

func RandomFilter(imgApi AbstractImageAPI, image []byte, contentType string, kernelSize, lower, higher int64, normalize bool) ([]byte, error) {
	return imgApi.ProcessImage(image, contentType, "randomFilter", util.RandomFilterQuery(kernelSize, lower, higher, normalize))
}

func InvertImage(imgApi AbstractImageAPI, image []byte, contentType string) ([]byte, error) {
	return imgApi.ProcessImage(image, contentType, "invert", "")

}

func SaturateImage(imgApi AbstractImageAPI, image []byte, contentType string, saturation int64) ([]byte, error) {
	saturationNorm := float32(saturation) / 100.0
	return imgApi.ProcessImage(image, contentType, "saturate", util.SaturateQuery(saturationNorm))

}
func EdgeDetect(imgApi AbstractImageAPI, image []byte, contentType string, lower, higher int64) ([]byte, error) {
	return imgApi.ProcessImage(image, contentType, "edgeDetection", util.EdgeDetectQuery(lower, higher))

}

func DilateImage(imgApi AbstractImageAPI, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	return imgApi.ProcessImage(image, contentType, "morphology", util.DilateQuery(kernelSize, iterations))
}

func ErodeImage(imgApi AbstractImageAPI, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	return imgApi.ProcessImage(image, contentType, "morphology", util.ErodeQuery(kernelSize, iterations))

}

func AddText(imgApi AbstractImageAPI, image []byte, contentType string, text string, fontScale float32, xPercentage float32, yPercentage float32) ([]byte, error) {
	return imgApi.ProcessImage(image, contentType, "text", util.AddTextQuery(text, fontScale, xPercentage, yPercentage))
}

func Reduced(imgApi AbstractImageAPI, image []byte, contentType string, quality float32) ([]byte, error) {

	return imgApi.ProcessImage(image, contentType, "reduced", util.ReduceQuery(quality))

}

func Shuffle(imgApi AbstractImageAPI, image []byte, contentType string, partitions int64) ([]byte, error) {

	return imgApi.ProcessImage(image, contentType, "shuffle", util.ShuffleQuery(partitions))
}
