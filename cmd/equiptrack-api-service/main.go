package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lukaslechman/equiptrack-webapi/api"
	"github.com/lukaslechman/equiptrack-webapi/internal/equiptrack"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("EQUIPTRACK_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("EQUIPTRACK_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	// request routings
	handleFunctions := &equiptrack.ApiHandleFunctions{
		EquipmentRegistryAPI: equiptrack.NewEquipmentRegistryApi(),
	}
	equiptrack.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)

}
