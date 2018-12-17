package database

import (
	"github.com/coc1961/go/crud/entities"
)

//Database Database
type Database interface {
	Get(id string) (*entities.Entity, error)
	Find(query map[string]interface{}) ([]*entities.Entity, error)
	Delete(id string) error
	Insert(entity *entities.Entity) error
	Update(id string, entity *entities.Entity) error
}
