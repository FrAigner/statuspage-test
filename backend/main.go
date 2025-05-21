package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialisiere die Datenbank
	initDB()

	r := gin.Default()

	// Debug-Modus aktivieren
	gin.SetMode(gin.DebugMode)

	// Logger aktivieren
	r.Use(gin.Logger())

	// Recovery aktivieren
	r.Use(gin.Recovery())

	// CORS-Konfiguration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Root-Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Statuspage API is running",
		})
	})

	// Health-Check Route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// API-Routen
	api := r.Group("/api")
	{
		// Service-Status Routen
		api.GET("/services", getServices)
		api.POST("/services", createService)
		api.PUT("/services/:id", updateService)
		api.DELETE("/services/:id", deleteService)

		// Incident Routen
		api.GET("/incidents", getIncidents)
		api.POST("/incidents", createIncident)
		api.PUT("/incidents/:id", updateIncident)
		api.GET("/incidents/:id", getIncidentDetails)

		// Component Routen
		api.GET("/components", getComponents)
		api.POST("/components", createComponent)
		api.PUT("/components/:id", updateComponent)
		api.DELETE("/components/:id", deleteComponent)

		// Tag Routen
		api.GET("/tags", getTags)
		api.POST("/tags", createTag)
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	log.Fatal(r.Run(":8080"))
}
