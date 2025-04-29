package util

import (
	"fmt"
	"net/url"
)

func SaturateQuery(saturationValue float32) string {
	return fmt.Sprintf("?saturation=%0.2f", saturationValue)
}

func EdgeDetectQuery(tLower, tHigher int64) string {
	return fmt.Sprintf("?lower=%d&higher=%d", tLower, tHigher)
}

func DilateQuery(kernelSize, iterations int64) string {
	return fmt.Sprintf("?kernelSize=%d&iterations=%d&type=Dilate", kernelSize, iterations)
}

func ErodeQuery(kernelSize, iterations int64) string {
	return fmt.Sprintf("?kernelSize=%d&iterations=%d&type=Erode", kernelSize, iterations)
}

func ReduceQuery(quality float32) string {
	return fmt.Sprintf("?quality=%0.2f", quality)
}

func AddTextQuery(text string, fontscale, xPerc, yPerc float32) string {

	uriText := url.QueryEscape(text)

	return fmt.Sprintf("?text=%s&fontScale=%0.2f&xPerc=%0.2f&yPerc=%0.2f", uriText, fontscale, xPerc, yPerc)
}

func RandomFilterQuery(kernelSize, minVal, maxVal int64, normalize bool) string {

	normalizeStr := "true"
	if !normalize {
		normalizeStr = "false"
	}

	return fmt.Sprintf("?minVal=%d&maxVal=%d&kernelSize=%d&normalize=%s", minVal, maxVal, kernelSize, normalizeStr)
}

func ShuffleQuery(partitions int64) string {
	return fmt.Sprintf("?partitions=%d", partitions)
}
