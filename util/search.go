package util

import (
	"fmt"
	"net/url"

	"github.com/raitonoberu/ytsearch"
)

func IsURL(input string) bool {

	u, err := url.ParseRequestURI(input)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""

}

func GetURLFromQuery(query string) (string, error) {
	search := ytsearch.VideoSearch(query)
	result, err := search.Next()

	if err != nil || len(result.Videos) == 0 {
		return "", fmt.Errorf("Query: \" %s \" did not produce any results", query)
	}

	url := result.Videos[0].URL
	return url, nil
}
