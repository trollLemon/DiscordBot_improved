package factories

import (
	"bot/Core/Services/ImageManip"
	"fmt"
)

const (
	Imagemanip Service = 0
)

const (

	ImageManipEndpoint = "http://image:8080/api"
)

func CreateApiService(service Service) (imagemanip.ImageAPI, error) {

	switch service {

	case Imagemanip:
		return imagemanip.NewImageAPIWrapper(ImageManipEndpoint), nil
	default:
		return nil, fmt.Errorf("Cannot create non-existant service")

	}

}
