package database

import (
	"github.com/coc1961/go/jsonutil"

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

	// Creo la query y busco por id
	q := d.collection.FindId(bson.ObjectIdHex(id))
	err = q.One(&tmp)
	if err != nil {
		return nil, err
	}

	// Convierto en json
	json := jsonutil.NewFromMap(&tmp)

	// creo la entidad
	entity, err1 := d.definition.New(json.JSON())
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
	err := d.collection.RemoveId(bson.ObjectIdHex(id))
	return err
}

// Insert insert
func (d *MongoDB) Insert(entity *entities.Entity) error {
	var err error
	var tmp = make(map[string]interface{})
	var oid bson.ObjectId

	// Genero el Id
	oid = bson.NewObjectId()

	_, err = d.definition.Validate(entity.JSON())
	if err != nil {
		return err
	}

	//entity.Add("_id").Set(oid)

	err = json.Unmarshal([]byte(entity.JSON()), &tmp)
	if err != nil {
		return err
	}

	b := make(bson.M, len(tmp))
	for k, v := range tmp {
		b[k] = v
	}
	b["_id"] = oid

	entity.Add("_id").Set(oid)

	err = d.collection.Insert(&b)
	if err != nil {
		return err
	}

	//	d.collection.Insert(bson.M{"_id": oid, "foo": "bar"})

	return nil
}

// Update update
func (d *MongoDB) Update(id string, entity *entities.Entity) error {
	//d.collection.Update()
	return nil
}
