package equiptrack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukaslechman/equiptrack-webapi/internal/db_service"
)

type equipmentUpdater = func(
	ctx *gin.Context,
	equipment *Equipment,
) (updatedEquipment *Equipment, responseContent interface{}, status int)

func updateEquipmentFunc(ctx *gin.Context, updater equipmentUpdater) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
			"error":   "db_service not found",
		})
		return
	}

	db, ok := value.(db_service.DbService[Equipment])
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of required type",
			"error":   "cannot cast db_service context to db_service.DbService",
		})
		return
	}

	equipmentId := ctx.Param("equipmentId")

	equipment, err := db.FindDocument(ctx, equipmentId)
	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Equipment not found",
			"error":   err.Error(),
		})
		return
	default:
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to load equipment from database",
			"error":   err.Error(),
		})
		return
	}

	updatedEquipment, responseObject, status := updater(ctx, equipment)

	if updatedEquipment != nil {
		err = db.UpdateDocument(ctx, equipmentId, updatedEquipment)
	} else {
		err = nil
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Equipment was deleted while processing the request",
			"error":   err.Error(),
		})
	default:
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to update equipment in database",
			"error":   err.Error(),
		})
	}
}
