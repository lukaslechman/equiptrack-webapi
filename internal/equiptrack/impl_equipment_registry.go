package equiptrack

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implEquipmentRegistryAPI struct {
}

func NewEquipmentRegistryApi() EquipmentRegistryAPI {
	return &implEquipmentRegistryAPI{}
}

func (o implEquipmentRegistryAPI) ListEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentRegistryAPI) CreateEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentRegistryAPI) GetEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentRegistryAPI) UpdateEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentRegistryAPI) UpdateEquipmentStatus(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentRegistryAPI) DeleteEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
