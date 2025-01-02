package database

import (
	"fmt"
)

type Repository struct {
	db DatabaseService
}

func NewRepository(db DatabaseService) *Repository {

	return &Repository{db: db}
}

func (r *Repository) Add(item string) error {

	if r.db.IsPresent(item) {
		return fmt.Errorf("Cannot add duplicate item %s", item)
	}
	err := r.db.Insert(item)
	return err
}

func (r *Repository) Remove(item string) error {

	err := r.db.Delete(item)

	return err

}

func (r *Repository) GetRandN(n int) ([]string, error) {

	data, err := r.db.FetchRandom(n)

	if err != nil {
		return nil, fmt.Errorf("Could not fetch data from database: %s", err.Error())
	}

	return data, nil

}

func (r *Repository) GetAllItems() ([]string, error) {

	items, err := r.db.GetAll()

	if err != nil {
		return nil, err
	}

	return items, nil
}
