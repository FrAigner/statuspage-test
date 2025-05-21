package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Service Handlers
func getServices(c *gin.Context) {
	var services []Service
	if err := db.Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func createService(c *gin.Context) {
	var service Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.ID = uuid.New()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	if err := db.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func updateService(c *gin.Context) {
	id := c.Param("id")
	var service Service

	if err := db.First(&service, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service nicht gefunden"})
		return
	}

	var updateData Service
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Name = updateData.Name
	service.Description = updateData.Description
	service.Status = updateData.Status
	service.UpdatedAt = time.Now()

	if err := db.Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

func deleteService(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Service{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Incident Handlers
func getIncidents(c *gin.Context) {
	var incidents []Incident
	if err := db.Preload("Service").Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func createIncident(c *gin.Context) {
	var incident Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incident.ID = uuid.New()
	incident.CreatedAt = time.Now()
	incident.UpdatedAt = time.Now()

	if err := db.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Lade die Service-Informationen
	if err := db.Preload("Service").First(&incident, "id = ?", incident.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func updateIncident(c *gin.Context) {
	id := c.Param("id")
	var incident Incident

	if err := db.First(&incident, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vorfall nicht gefunden"})
		return
	}

	var updateData Incident
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incident.Title = updateData.Title
	incident.Description = updateData.Description
	incident.Status = updateData.Status
	incident.Impact = updateData.Impact
	incident.UpdatedAt = time.Now()

	if updateData.Status == "resolved" && incident.ResolvedAt == nil {
		now := time.Now()
		incident.ResolvedAt = &now
	}

	if err := db.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func getIncidentDetails(c *gin.Context) {
	id := c.Param("id")
	var incident Incident

	if err := db.Preload("Service").First(&incident, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vorfall nicht gefunden"})
		return
	}

	var updates []IncidentUpdate
	if err := db.Where("incident_id = ?", id).Find(&updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incident": incident,
		"updates":  updates,
	})
}
