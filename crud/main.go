package main

import (
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/crudframework"
	"github.com/coc1961/go/crud/entities"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = "/home/carlos/gopath/src/github.com/coc1961/go/crud"

	cf := crudframework.New(dir)
	cf.Load()

	cf.AddHandler("prueba", &Test{})

}

// Test test
type Test struct {
}

// OnAfterInsert OnAfterInsert
func (t *Test) OnAfterInsert(entity *entities.Entity) error {
	return nil
}

// OnBeforeInsert OnBeforeInsert
func (t *Test) OnBeforeInsert(entity *entities.Entity) error {
	return nil
}

// OnAfterUpdate OnAfterUpdate
func (t *Test) OnAfterUpdate(entity *entities.Entity, actualEntity *entities.Entity) error {
	return nil
}

// OnBeforeUpdate OnBeforeUpdate
func (t *Test) OnBeforeUpdate(entity *entities.Entity, actualEntity *entities.Entity) error {
	return nil
}

// OnAfterDelete OnAfterDelete
func (t *Test) OnAfterDelete(entity *entities.Entity) error {
	return nil
}

// OnBeforeDelete OnBeforeDelete
func (t *Test) OnBeforeDelete(entity *entities.Entity) error {
	return nil
}
