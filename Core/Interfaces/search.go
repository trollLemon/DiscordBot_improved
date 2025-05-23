package Interfaces

type Search interface {
	SearchWithQuery(query string) (string, error)
}
