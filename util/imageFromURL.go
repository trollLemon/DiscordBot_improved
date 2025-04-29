package util

import (
	"io"
	"net/http"
)

func GetImageFromURL(imageURL string) ([]byte, string, error) {

	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, "", err
	}

	format := resp.Header.Get("Content-Type")

	return bytes, format, nil

}
