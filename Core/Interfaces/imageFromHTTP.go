package Interfaces

type ImageFromHTTP interface {
	urlToBytes(url string) ([]byte, error)
}
