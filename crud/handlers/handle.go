package handlers

import (
	"net/http"

	"github.com/coc1961/go/crud/entities"
	"github.com/gin-gonic/gin"
)

// Handle entity handle
type Handle struct {
	name          string
	eventHandlers []EventHandler
	definition    *entities.Definition
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
func New(definition *entities.Definition) *Handle {
	return &Handle{name: definition.Name(), eventHandlers: make([]EventHandler, 0), definition: definition}
}

// AddEventHandler add handle to entity
func (h *Handle) AddEventHandler(handle EventHandler) error {
	h.eventHandlers = append(h.eventHandlers, handle)
	return nil
}

//Register Register service
func (h *Handle) Register(basePath string, router *gin.Engine) {
	router.GET(basePath+"/"+h.definition.Name(), h.Get)
	router.GET(basePath+"/"+h.definition.Name()+"/:id", h.Get)
	router.POST(basePath+"/"+h.definition.Name(), h.Post)
	router.PUT(basePath+"/"+h.definition.Name()+"/:id", h.Put)
	router.DELETE(basePath+"/"+h.definition.Name()+"/:id", h.Delete)
}

// Ver https://github.com/gin-gonic/gin#using-get-post-put-patch-delete-and-options

// Get get
func (h *Handle) Get(c *gin.Context) {
	c.String(http.StatusOK, "Not Implemented")
}

// Post get
func (h *Handle) Post(c *gin.Context) {
	c.String(http.StatusOK, "Not Implemented")
}

// Put get
func (h *Handle) Put(c *gin.Context) {
	c.String(http.StatusOK, "Not Implemented")
}

// Delete get
func (h *Handle) Delete(c *gin.Context) {
	c.String(http.StatusOK, "Not Implemented")
}
