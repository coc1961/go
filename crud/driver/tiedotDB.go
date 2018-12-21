package driver

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/coc1961/go/config"
	"github.com/coc1961/go/jsonutil"
	schema "github.com/lestrrat/go-jsschema"

	"encoding/json"

	"github.com/coc1961/go/crud/database"
	"github.com/coc1961/go/crud/entities"

	"github.com/HouzuoGuo/tiedot/db"
)

/**********
* Factory
**********/

// TiedotDBFactory TiedotDBFactory
type TiedotDBFactory struct {
}

// Create Return a TiedotDB Instance
func (f *TiedotDBFactory) Create() (database.Database, error) {
	return NewTiedot(config.Get().DatabaseConfig[0], config.Get().DatabaseConfig[1])
}

/**********************
* Mongo Db Database
***********************/

// TiedotDB implementacion en mongo
type TiedotDB struct {
	definition *entities.Definition
	db         *db.DB
	database   string
}

// NewTiedot nueva instancia
func NewTiedot(path, database string) (database.Database, error) {

	myDBDir := filepath.Join(path, database)

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		return nil, err
	}

	return &TiedotDB{definition: nil, database: myDBDir, db: myDB}, nil
}

func tiedottiedotRecoverError(reco interface{}) error {
	if reco != nil {
		return fmt.Errorf("%s", reco)
	}
	return nil
}

// SetDefinition Set Definition
func (d *TiedotDB) SetDefinition(definition *entities.Definition) {
	d.definition = definition
	if err := d.db.Create(definition.Name()); err != nil {
		log.Print("Collection " + definition.Name() + " exists")
	}
}

func (d *TiedotDB) getCollection() *db.Col {
	col := d.db.Use(d.definition.Name())
	prop := d.definition.Schema().Properties
	for k := range prop {
		col.Index([]string{k})
	}

	return col
}

func (d *TiedotDB) convertID(id string) (int, error) {
	return strconv.Atoi(id)
}

// Get get
func (d *TiedotDB) Get(id string) (*entities.Entity, error) {
	var err error
	var entity *entities.Entity
	var tmp = make(map[string]interface{})
	var nid int

	db := d.getCollection()

	nid, err = d.convertID(id)
	if err != nil {
		return nil, err
	}

	tmp, err = db.Read(nid)
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
func (d *TiedotDB) Find(query map[string]interface{}) ([]*entities.Entity, error) {

	dbc := d.getCollection()

	str := "["
	coma := ""

	for s := range query {
		field := d.definition.Schema().Properties[s]
		if field == nil {
			return nil, errors.New("attribute " + s + " not exists")
		}
		typ := field.Type[0]
		if typ == schema.IntegerType || typ == schema.NumberType {
			str = str + coma + "{\"eq\": " + query[s].(string) + ","
		} else {
			str = str + coma + "{\"eq\": \"" + query[s].(string) + "\","
		}
		str = str + "\"in\": " + "[\"" + s + "\"]"
		str = str + "}"
		coma = ","
	}

	str = str + "]"
	var q interface{}
	json.Unmarshal([]byte(str), &q)

	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

	if err := db.EvalQuery(q, dbc, &queryResult); err != nil {
		return nil, err
	}
	ret := make([]*entities.Entity, 0)

	for id := range queryResult {
		// To get query result document, simply read it
		tmp, err := d.Get(strconv.Itoa(id))
		if err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}

	/*
		// Query result are document IDs
		for id := range queryResult {
			// To get query result document, simply read it
			readBack, err := feeds.Read(id)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Query returned document %v\n", readBack)
		}
	*/
	/*
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
	*/
	return ret, nil
}

// Delete delete
func (d *TiedotDB) Delete(id string) error {
	db := d.getCollection()
	nid, err := d.convertID(id)
	if err != nil {
		return err
	}
	err = db.Delete(nid)
	return err
}

// Insert insert
func (d *TiedotDB) Insert(entity *entities.Entity) error {
	var err error
	var tmp = make(map[string]interface{})
	var nid int

	_, err = d.definition.Validate(entity.JSON())
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(entity.JSON()), &tmp)
	if err != nil {
		return err
	}

	db := d.getCollection()

	nid, err = db.Insert(tmp)
	if err != nil {
		return err
	}
	tmp["_id"] = strconv.Itoa(nid)
	db.Update(nid, tmp)

	entity.Add("_id").Set(strconv.Itoa(nid))

	return err
}

// Update update
func (d *TiedotDB) Update(id string, entity *entities.Entity) error {
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

	db := d.getCollection()
	nid, _ := d.convertID(id)
	tmp["_id"] = strconv.Itoa(nid)
	entity.Add("_id").Set(tmp["_id"])

	err = db.Update(nid, tmp)

	return err
}
