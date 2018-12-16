package database

import (
	"github.com/coc1961/go/crud/entities"
)

//Database Database
type Database interface {
	Get(id string) (*entities.Entity, error)
	Find(query map[string]string) ([]*entities.Entity, error)
	Delete(id string) error
	Insert(id string, entity *entities.Entity) error
	Update(id string, entity *entities.Entity) error
}
