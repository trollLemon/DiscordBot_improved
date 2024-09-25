package util

import (
	"github.com/raitonoberu/ytsearch"
	"net/url"
)

func IsURL(input string) bool {

	u, err := url.ParseRequestURI(input)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""

}

func GetURLFromQuery(query string) (string, error) {
	search := ytsearch.VideoSearch(query)
	result, err := search.Next()
	url := result.Videos[0].URL
	return url,err
}
