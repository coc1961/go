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

	var factory database.Factory

	factory = &database.MongoDBFactory{}

	cf.Load(factory)

	cf.AddEventHandler("prueba", &TestEventHandler{})

	cf.Start()
}

// TestEventHandler test
type TestEventHandler struct {
}

// OnAfterInsert OnAfterInsert
func (t *TestEventHandler) OnAfterInsert(db database.Database, entity *entities.Entity, err error) error {
	return nil
}

// OnBeforeInsert OnBeforeInsert
func (t *TestEventHandler) OnBeforeInsert(db database.Database, entity *entities.Entity) error {
	return nil
}

// OnAfterUpdate OnAfterUpdate
func (t *TestEventHandler) OnAfterUpdate(db database.Database, entity *entities.Entity, actualEntity *entities.Entity, err error) error {
	return nil
}

// OnBeforeUpdate OnBeforeUpdate
func (t *TestEventHandler) OnBeforeUpdate(db database.Database, entity *entities.Entity, actualEntity *entities.Entity) error {
	return nil
}

// OnAfterDelete OnAfterDelete
func (t *TestEventHandler) OnAfterDelete(db database.Database, entity *entities.Entity, err error) error {
	return nil
}

// OnBeforeDelete OnBeforeDelete
func (t *TestEventHandler) OnBeforeDelete(db database.Database, entity *entities.Entity) error {
	return nil
}
