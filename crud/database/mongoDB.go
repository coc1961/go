package database

import (
	"encoding/json"

	"github.com/coc1961/go/crud/entities"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB implementacion en mongo
type MongoDB struct {
	definition *entities.Definition
	collection *mgo.Collection
}

// NewMongo nueva instancia
func NewMongo(ip, database string, definition *entities.Definition) (Database, error) {
	session, err := mgo.Dial(ip)
	if err != nil {
		return nil, err
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(database).C(definition.Name())
	return &MongoDB{collection: c, definition: definition}, nil
}

// Get get
func (d *MongoDB) Get(id string) (*entities.Entity, error) {
	var err error
	var entity *entities.Entity
	var tmp = make(map[string]interface{})
	var b []byte

	// Creo la query y busco por id
	q := d.collection.FindId(bson.ObjectIdHex(id))
	err = q.One(&tmp)
	if err != nil {
		return nil, err
	}

	// Convierto en json
	b, err = json.Marshal(tmp)
	if err != nil {
		return nil, err
	}

	// creo la entidad
	entity, err1 := d.definition.New(string(b))
	if err1 != nil {
		return nil, err
	}

	return entity, err
}

// Find find
func (d *MongoDB) Find(query map[string]string) ([]*entities.Entity, error) {
	return nil, nil
}

// Delete delete
func (d *MongoDB) Delete(id string) error {
	return nil
}

// Insert insert
func (d *MongoDB) Insert(id string, entity *entities.Entity) error {
	return nil
}

// Update update
func (d *MongoDB) Update(id string, entity *entities.Entity) error {
	return nil
}
