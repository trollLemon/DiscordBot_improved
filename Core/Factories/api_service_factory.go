package factories

import (
	"bot/Core/Services/ImageManip"
	"fmt"
)

const (
	Imagemanip         Service = 4
	ImageManipEndpoint         = "http://image:8080/api/image"
)

func CreateImageAPIService(service Service) (*imagemanip.ImageAPIWrapper, error) {

	switch service {

	case Imagemanip:
		return imagemanip.NewImageAPIWrapper(ImageManipEndpoint), nil
	default:
		return nil, fmt.Errorf("invalid image manipulation service")

	}

}
