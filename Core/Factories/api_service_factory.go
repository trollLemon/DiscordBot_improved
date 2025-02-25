package factories

import (
	"bot/Core/Services"
	"bot/Core/Services/ImageManip"
	"fmt"
)

const (
	Imagemanip Service = 0
)

func CreateApiService(service Service, endpoint string) (services.IService, error) {

	switch service {

	case Imagemanip:
		return imagemanip.NewImageAPIWrapper(endpoint), nil
	default:
		return nil, fmt.Errorf("Cannot create non-existant service")

	}

}
