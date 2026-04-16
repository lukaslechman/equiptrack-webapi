package equiptrack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukaslechman/equiptrack-webapi/internal/db_service"
	"go.mongodb.org/mongo-driver/bson"
)

type implEquipmentRegistryAPI struct {
}

func NewEquipmentRegistryApi() EquipmentRegistryAPI {
	return &implEquipmentRegistryAPI{}
}

// GET /equipment
func (o *implEquipmentRegistryAPI) ListEquipment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error", "message": "db_service not found",
		})
		return
	}
	db, ok := value.(db_service.DbService[Equipment])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error", "message": "db_service cast failed",
		})
		return
	}

	// zostavaj filter z query params
	filter := bson.D{}

	if status := c.Query("status"); status != "" {
		filter = append(filter, bson.E{Key: "status", Value: status})
	}

	if department := c.Query("department"); department != "" {
		filter = append(filter, bson.E{Key: "location.department", Value: department})
	}

	if category := c.Query("category"); category != "" {
		filter = append(filter, bson.E{Key: "category", Value: category})
	}

	equipment, err := db.FindDocuments(c, filter)
	switch err {
	case nil:
		// vráť prázdny zoznam namiesto null
		if equipment == nil {
			equipment = []*Equipment{}
		}
		c.JSON(http.StatusOK, equipment)
	default:
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "Bad Gateway", "message": "Failed to load equipment", "error": err.Error(),
		})
	}
}

// POST /equipment
func (o *implEquipmentRegistryAPI) CreateEquipment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error", "message": "db_service not found",
		})
		return
	}
	db, ok := value.(db_service.DbService[Equipment])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error", "message": "db_service cast failed",
		})
		return
	}

	var equipment Equipment
	if err := c.BindJSON(&equipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request", "message": "Invalid request body", "error": err.Error(),
		})
		return
	}

	if equipment.Id == "" {
		equipment.Id = uuid.New().String()
	}

	if equipment.Status == "" {
		equipment.Status = "active"
	}

	err := db.CreateDocument(c, equipment.Id, &equipment)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, equipment)
	case db_service.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{
			"status": "Conflict", "message": "Equipment already exists", "error": err.Error(),
		})
	default:
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "Bad Gateway", "message": "Failed to create equipment", "error": err.Error(),
		})
	}
}

// GET /equipment/:equipmentId
func (o *implEquipmentRegistryAPI) GetEquipment(c *gin.Context) {
	updateEquipmentFunc(c, func(c *gin.Context, equipment *Equipment) (*Equipment, interface{}, int) {
		return nil, equipment, http.StatusOK
	})
}

// PUT /equipment/:equipmentId
func (o *implEquipmentRegistryAPI) UpdateEquipment(c *gin.Context) {
	updateEquipmentFunc(c, func(c *gin.Context, equipment *Equipment) (*Equipment, interface{}, int) {
		var updatedData Equipment
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			return nil, gin.H{
				"status": http.StatusBadRequest, "message": "Invalid request body", "error": err.Error(),
			}, http.StatusBadRequest
		}

		updatedData.Id = equipment.Id

		return &updatedData, updatedData, http.StatusOK
	})
}

// PATCH /equipment/:equipmentId  (zmena statusu)
func (o *implEquipmentRegistryAPI) UpdateEquipmentStatus(c *gin.Context) {
	updateEquipmentFunc(c, func(c *gin.Context, equipment *Equipment) (*Equipment, interface{}, int) {
		var body EquipmentStatusUpdate
		if err := c.ShouldBindJSON(&body); err != nil {
			return nil, gin.H{
				"status": http.StatusBadRequest, "message": "Invalid request body", "error": err.Error(),
			}, http.StatusBadRequest
		}

		equipment.Status = body.Status
		return equipment, equipment, http.StatusOK
	})
}

// DELETE /equipment/:equipmentId
func (o *implEquipmentRegistryAPI) DeleteEquipment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error", "message": "db_service not found",
		})
		return
	}
	db, ok := value.(db_service.DbService[Equipment])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error", "message": "db_service cast failed",
		})
		return
	}

	equipmentId := c.Param("equipmentId")
	err := db.DeleteDocument(c, equipmentId)
	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Not Found", "message": "Equipment not found", "error": err.Error(),
		})
	default:
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "Bad Gateway", "message": "Failed to delete equipment", "error": err.Error(),
		})
	}
}
