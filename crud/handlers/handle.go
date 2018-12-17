package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/coc1961/go/crud/database"
	"github.com/coc1961/go/crud/entities"
	"github.com/coc1961/go/jsonutil"
	"github.com/gin-gonic/gin"
)

// Handle entity handle
type Handle struct {
	name          string
	eventHandlers []EventHandler
	definition    *entities.Definition
	database      database.Database
}

/*********
* Handlers and Validators
**/

// EventHandler Crud Handler
type EventHandler interface {
	OnAfterInsert(entity *entities.Entity, err error) error
	OnBeforeInsert(entity *entities.Entity) error

	OnAfterUpdate(entity *entities.Entity, actualEntity *entities.Entity, err error) error
	OnBeforeUpdate(entity *entities.Entity, actualEntity *entities.Entity) error

	OnAfterDelete(entity *entities.Entity, err error) error
	OnBeforeDelete(entity *entities.Entity) error
}

// New nuevo handler
func New(definition *entities.Definition, db database.Database) *Handle {
	return &Handle{name: definition.Name(), eventHandlers: make([]EventHandler, 0), definition: definition, database: db}
}

// AddEventHandler add handle to entity
func (h *Handle) AddEventHandler(handle EventHandler) error {
	h.eventHandlers = append(h.eventHandlers, handle)
	return nil
}

//Register Register service
func (h *Handle) Register(basePath string, router *gin.Engine) {
	router.GET(basePath+"/"+h.definition.Name(), h.Find)
	router.GET(basePath+"/"+h.definition.Name()+"/:id", h.Get)
	router.POST(basePath+"/"+h.definition.Name(), h.Post)
	router.PUT(basePath+"/"+h.definition.Name()+"/:id", h.Put)
	router.DELETE(basePath+"/"+h.definition.Name()+"/:id", h.Delete)
}

// Ver https://github.com/gin-gonic/gin#using-get-post-put-patch-delete-and-options

// Find find
func (h *Handle) Find(c *gin.Context) {
	values := c.Request.URL.Query()
	result := make(map[string]interface{})
	for k, v := range values {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}

	ent, err := h.database.Find(result)
	if err == nil {
		json := "["
		coma := ""
		for _, e := range ent {
			json = json + coma + e.JSON()
			coma = ","
		}
		json = json + "]"
		c.String(http.StatusOK, json)
	} else {
		c.String(http.StatusInternalServerError, errorToJSON(err))
	}
}

// Get get
func (h *Handle) Get(c *gin.Context) {
	id := c.Param("id")
	ent, err := h.database.Get(id)
	if err == nil {
		c.String(http.StatusOK, ent.JSON())
	} else {
		c.String(http.StatusInternalServerError, errorToJSON(err))
	}
}

// Post get
func (h *Handle) Post(c *gin.Context) {
	var err error
	txt, _ := ioutil.ReadAll(c.Request.Body)
	ent, _ := h.definition.New(string(txt))
	for _, ev := range h.eventHandlers {
		err1 := ev.OnBeforeInsert(ent)
		if err1 != nil {
			err = err1
			break
		}
	}
	err = h.database.Insert(ent)

	for _, ev := range h.eventHandlers {
		err1 := ev.OnAfterInsert(ent, err)
		if err1 != nil {
			err = err1
			break
		}
	}
	if err != nil {
		c.String(http.StatusBadRequest, errorToJSON(err))
	} else {
		c.String(http.StatusOK, ent.JSON())
	}
}

// Put get
func (h *Handle) Put(c *gin.Context) {
	txt, _ := ioutil.ReadAll(c.Request.Body)
	ent, _ := h.definition.New(string(txt))

	id := c.Param("id")
	entAnt, err := h.database.Get(id)

	for _, ev := range h.eventHandlers {
		err1 := ev.OnBeforeUpdate(ent, entAnt)
		if err1 != nil {
			err = err1
			break
		}
	}

	if entAnt != nil {
		err = h.database.Update(id, ent)
	}

	for _, ev := range h.eventHandlers {
		err1 := ev.OnAfterUpdate(ent, entAnt, err)
		if err1 != nil {
			err = err1
			break
		}
	}

	if err != nil {
		c.String(http.StatusBadRequest, errorToJSON(err))
	} else {
		c.String(http.StatusOK, ent.JSON())
	}
}

// Delete get
func (h *Handle) Delete(c *gin.Context) {
	id := c.Param("id")
	ent, err := h.database.Get(id)

	for _, ev := range h.eventHandlers {
		err1 := ev.OnBeforeDelete(ent)
		if err1 != nil {
			err = err1
			break
		}
	}

	err = h.database.Delete(id)

	for _, ev := range h.eventHandlers {
		err1 := ev.OnAfterDelete(ent, err)
		if err1 != nil {
			err = err1
			break
		}
	}
	if err != nil {
		c.String(http.StatusBadRequest, errorToJSON(err))
	} else {
		c.String(http.StatusOK, ent.JSON())
	}
}

// Convierto un error en json
func errorToJSON(err error) string {
	json := jsonutil.New()
	json.Add("error").Set(err.Error())
	return json.JSON()
}
