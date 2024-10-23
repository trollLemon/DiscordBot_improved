package database

import (
)



/* DatabaseService
 * 
 * Defines an interface for interacting with a database, 
 * + adding data 
 * + removing data
 * + getting n random items from the database 
 *   - This is required for the random play functionality
 */
type DatabaseService interface {
	
	Insert(item string) error   // Insert into database
	Delete(item string) error   // Remove from database
	FetchRandom(n int) ([]string,  error)   // Get n random datapoints from database
	IsPresent(item string) bool // check is item is in database
	GetAll() ([]string, error)
}




