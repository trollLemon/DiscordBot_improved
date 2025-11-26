package store

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

var (
	ErrItemNotFound  = errors.New("the item is not in the database")
	ErrDuplicateItem = errors.New("the item is already in the database")
	ErrEmpty         = errors.New("the database is empty")
	ErrGeneral       = errors.New("error communicating with database")
)

type RandomWords struct {
	store *Redis
}

func errorChecker(err error) error {

	if errors.Is(err, errNotFound) {
		return fmt.Errorf("could not complete action because %w", ErrItemNotFound)
	}
	if errors.Is(err, errDuplicate) {
		return fmt.Errorf("could not complete action because %w", ErrDuplicateItem)
	}
	if errors.Is(err, errEmpty) {
		return fmt.Errorf("could not complete action because %w", ErrEmpty)
	}
	if err != nil {
		return ErrGeneral
	}

	return nil
}

func NewRandomWords(store *Redis) *RandomWords {
	return &RandomWords{
		store: store,
	}
}

func (r *RandomWords) Insert(item string) error {

	insertErr := r.store.Insert(item)

	if insertErr != nil {
		log.Err(insertErr).Msg("database action failed")
		return errorChecker(insertErr)
	}

	return nil

}

func (r *RandomWords) Delete(item string) error {
	deleteErr := r.store.Delete(item)

	if deleteErr != nil {
		log.Err(deleteErr).Msg("database action failed")
		return errorChecker(deleteErr)
	}

	return nil
}

func (r *RandomWords) GetRandom(n int) ([]string, error) {
	items, err := r.store.FetchRandom(n)

	if err != nil {
		log.Err(err).Msg("database action failed")
		return nil, errorChecker(err)
	}

	return items, nil
}

func (r *RandomWords) GetAll() ([]string, error) {
	items, err := r.store.GetAll()

	if err != nil {
		log.Err(err).Msg("database action failed")
		return nil, errorChecker(err)

	}

	return items, nil
}
