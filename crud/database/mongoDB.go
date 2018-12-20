package database

import (
	"errors"
	"fmt"

	"github.com/coc1961/go/config"
	"github.com/coc1961/go/jsonutil"

	"encoding/json"

	"github.com/coc1961/go/crud/entities"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**********
* Factory
**********/

// MongoDBFactory MongoDBFactory
type MongoDBFactory struct {
}

// Create Return a Mongodb Instance
func (f *MongoDBFactory) Create() (Database, error) {
	return NewMongo(config.Get().DatabaseConfig[0], config.Get().DatabaseConfig[1])
}

/**********************
* Mongo Db Database
***********************/

// MongoDB implementacion en mongo
type MongoDB struct {
	definition *entities.Definition
	collection *mgo.Collection
	session    *mgo.Session
	ip         string
	database   string
}

// NewMongo nueva instancia
func NewMongo(ip, database string) (Database, error) {
	session, err := mgo.Dial(ip)
	if err != nil {
		return nil, err
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return &MongoDB{collection: nil, definition: nil, session: session, ip: ip, database: database}, nil
}

func recoverError(reco interface{}) error {
	if reco != nil {
		return fmt.Errorf("%s", reco)
	}
	return nil
}

func toObjectID(id string) (oid bson.ObjectId, err error) {
	defer func() {
		err = recoverError(recover())
	}()
	oid = bson.ObjectIdHex(id)
	err = nil
	return oid, err
}

// SetDefinition Set Definition
func (d *MongoDB) SetDefinition(definition *entities.Definition) {
	c := d.session.DB(d.database).C(definition.Name())
	d.collection = c
	d.definition = definition
}

// Get get
func (d *MongoDB) Get(id string) (*entities.Entity, error) {
	var err error
	var entity *entities.Entity
	var tmp = make(map[string]interface{})

	// Creo la query y busco por id
	oid, err := toObjectID(id)
	if err != nil {
		return nil, err
	}
	q := d.collection.FindId(oid)
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
func (d *MongoDB) Find(query map[string]interface{}) ([]*entities.Entity, error) {
	b := make(bson.M, len(query))
	for k, v := range query {
		b[k] = v
	}
	var result []interface{}
	iter := d.collection.Find(&b).Limit(config.Get().FindLimit).Iter()
	err := iter.All(&result)

	ret := make([]*entities.Entity, 0)
	if result != nil {
		for _, v := range result {
			tmp := v.(bson.M)
			by, err1 := json.Marshal(&tmp)
			if err1 == nil {
				ent, err1 := d.definition.New(string(by))
				if err1 == nil {
					ret = append(ret, ent)
				}

			}
		}
	}

	return ret, err
}

// Delete delete
func (d *MongoDB) Delete(id string) error {
	oid, err := toObjectID(id)
	if err != nil {
		return err
	}
	err = d.collection.RemoveId(oid)
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

	return err
}

// Update update
func (d *MongoDB) Update(id string, entity *entities.Entity) error {
	var tmp = make(map[string]interface{})

	// Busco por id
	entAnt, err := d.Get(id)
	if err != nil {
		return err
	}

	if entAnt == nil {
		err = errors.New("Not Found")
		return err
	}

	_, err = d.definition.Validate(entity.JSON())
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(entity.JSON()), &tmp)
	if err != nil {
		return err
	}

	b := make(bson.M, len(tmp))
	for k, v := range tmp {
		b[k] = v
	}
	oid, err := toObjectID(id)
	if err != nil {
		return err
	}
	b["_id"] = oid
	entity.Add("_id").Set(b["_id"])

	err = d.collection.UpdateId(b["_id"], &b)

	return err
}
