package database

type AbstractDatabaseService interface {
	Insert(item string) error            // Insert into database
	Delete(item string) error            // Remove from database
	FetchRandom(n int) ([]string, error) // Get n random datapoints from database
	GetAll() ([]string, error)
}
