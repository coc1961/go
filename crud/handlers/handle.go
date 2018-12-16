package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/coc1961/go/crud/database"
	"github.com/coc1961/go/crud/entities"
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
	OnAfterInsert(entity *entities.Entity) error
	OnBeforeInsert(entity *entities.Entity) error

	OnAfterUpdate(entity *entities.Entity, actualEntity *entities.Entity) error
	OnBeforeUpdate(entity *entities.Entity, actualEntity *entities.Entity) error

	OnAfterDelete(entity *entities.Entity) error
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
	c.String(http.StatusOK, "Not Implemented")
}

// Get get
func (h *Handle) Get(c *gin.Context) {
	id := c.Param("id")
	ent, err := h.database.Get(id)
	if err == nil {
		c.String(http.StatusOK, ent.JSON())
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

// Post get
func (h *Handle) Post(c *gin.Context) {
	var err error = nil
	txt, _ := ioutil.ReadAll(c.Request.Body)
	ent, _ := h.definition.New(string(txt))
	for _, ev := range h.eventHandlers {
		err := ev.OnBeforeInsert(ent)
		if err != nil {
			break
		}
	}
	err = h.database.Insert(ent)
	if err == nil {
		for _, ev := range h.eventHandlers {
			err = ev.OnAfterInsert(ent)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, ent.JSON())
	}
}

// Put get
func (h *Handle) Put(c *gin.Context) {
	for _, ev := range h.eventHandlers {
		txt, _ := ioutil.ReadAll(c.Request.Body)
		ent, _ := h.definition.New(string(txt))
		ev.OnBeforeUpdate(ent, nil)
	}
	for _, ev := range h.eventHandlers {
		txt, _ := ioutil.ReadAll(c.Request.Body)
		ent, _ := h.definition.New(string(txt))
		ev.OnAfterUpdate(ent, nil)
	}
	c.String(http.StatusOK, "Not Implemented")
}

// Delete get
func (h *Handle) Delete(c *gin.Context) {
	for _, ev := range h.eventHandlers {
		ev.OnBeforeDelete(nil)
	}
	for _, ev := range h.eventHandlers {
		ev.OnAfterDelete(nil)
	}
	c.String(http.StatusOK, "Not Implemented")
}
