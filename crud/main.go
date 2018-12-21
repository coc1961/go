package main

import (
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/crudframework"
	"github.com/coc1961/go/crud/database"
	"github.com/coc1961/go/crud/driver"
	"github.com/coc1961/go/crud/entities"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = "/home/carlos/gopath/src/github.com/coc1961/go/crud"

	// Creo el framework
	cf := crudframework.New(dir)

	// Configuro el driver de la base de datos
	var factory database.Factory
	//factory = &driver.MongoDBFactory{}
	factory = &driver.TiedotDBFactory{}

	// Inicializo el framework con el drive de base de datos
	cf.InitFramework(factory)

	// Agrego un validador suscripto a los eventos
	cf.AddEventHandler("prueba", &TestEventHandler{})

	// Inicio el Server
	cf.Start("8080")
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
	/*
		query := make(map[string]interface{}, 1)
		query["id"] = entity.Get("id").Value().(string)
		lst, err := db.Find(query)
		if err == nil {
			if lst != nil && len(lst) > 0 {
				err = errors.New("record duplicated")
			}
		}
		return err
	*/
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
