package util

import (
	"fmt"
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

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("cannot download image from url: %s. Status Code %d", imageURL, resp.StatusCode)
	}

	format := resp.Header.Get("Content-Type")

	return bytes, format, nil

}
