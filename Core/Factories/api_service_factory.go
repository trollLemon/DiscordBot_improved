package factories

import (
	"bot/Core/Services/ImageManip"
	"fmt"
	"net/http"
)

const (
	GoManip            Service = 4
	ImageManipEndpoint         = "http://image:8080/api/image"
)

func CreateImageAPIService(service Service) (imagemanip.AbstractImageAPI, error) {

	switch service {

	case GoManip:
		return imagemanip.NewGoManip(&http.Client{}, ImageManipEndpoint), nil

	default:
		return nil, fmt.Errorf("invalid image manipulation service")

	}

}
