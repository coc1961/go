package main

import (
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/crudframework"
	"github.com/coc1961/go/crud/database"
	"github.com/coc1961/go/crud/entities"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = "/home/carlos/gopath/src/github.com/coc1961/go/crud"

	cf := crudframework.New(dir)
	cf.Load(&database.MongoDBFactory{})

	cf.AddEventHandler("prueba", &TestEventHandler{})

	cf.Start()
}

// TestEventHandler test
type TestEventHandler struct {
}

// OnAfterInsert OnAfterInsert
func (t *TestEventHandler) OnAfterInsert(entity *entities.Entity, err error) error {
	return nil
}

// OnBeforeInsert OnBeforeInsert
func (t *TestEventHandler) OnBeforeInsert(entity *entities.Entity) error {
	return nil
}

// OnAfterUpdate OnAfterUpdate
func (t *TestEventHandler) OnAfterUpdate(entity *entities.Entity, actualEntity *entities.Entity, err error) error {
	return nil
}

// OnBeforeUpdate OnBeforeUpdate
func (t *TestEventHandler) OnBeforeUpdate(entity *entities.Entity, actualEntity *entities.Entity) error {
	return nil
}

// OnAfterDelete OnAfterDelete
func (t *TestEventHandler) OnAfterDelete(entity *entities.Entity, err error) error {
	return nil
}

// OnBeforeDelete OnBeforeDelete
func (t *TestEventHandler) OnBeforeDelete(entity *entities.Entity) error {
	return nil
}
