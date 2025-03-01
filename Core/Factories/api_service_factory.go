package factories

import (
	"bot/Core/Services/ImageManip"
	"fmt"
)

const (

	Imagemanip Service = 4
	ImageManipEndpoint = "http://image:8080/api"
)

func CreateImageAPIService(service Service) (imagemanip.ImageAPI, error) {
	
	switch service {

	case Imagemanip:
		return imagemanip.NewImageAPIWrapper(ImageManipEndpoint), nil
	default:
		return nil, fmt.Errorf("Invalid image manipulation service")

	}

}
