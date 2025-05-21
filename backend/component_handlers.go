package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Tag Handlers
func getTags(c *gin.Context) {
	var tags []Tag
	if err := db.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func createTag(c *gin.Context) {
	var tag Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag.ID = uuid.New().String()
	if err := db.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// Component Handlers
func getComponents(c *gin.Context) {
	var components []Component

	// Filter by tags if provided
	tags := c.QueryArray("tags")

	query := db.Preload("Tags")

	if len(tags) > 0 {
		query = query.Joins("JOIN component_tags ON components.id = component_tags.component_id").
			Joins("JOIN tags ON component_tags.tag_id = tags.id").
			Where("tags.name IN (?)", tags).
			Group("components.id")
	}

	if err := query.Find(&components).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, components)
}

func createComponent(c *gin.Context) {
	var input struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Status      string   `json:"status"`
		TagNames    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	component := Component{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	// Begin Transaction
	tx := db.Begin()

	if err := tx.Create(&component).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle Tags
	for _, tagName := range input.TagNames {
		var tag Tag

		// Try to find existing tag
		err := tx.Where("name = ?", tagName).First(&tag).Error
		if err != nil {
			// Create new tag if it doesn't exist
			tag = Tag{
				ID:   uuid.New().String(),
				Name: tagName,
			}
			if err := tx.Create(&tag).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Add tag to component
		if err := tx.Create(&ComponentTag{ComponentID: component.ID, TagID: tag.ID}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	// Fetch the complete component with tags
	db.Preload("Tags").First(&component, "id = ?", component.ID)

	c.JSON(http.StatusCreated, component)
}

func updateComponent(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Status      string   `json:"status"`
		TagNames    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.Begin()

	var component Component
	if err := tx.Preload("Tags").First(&component, "id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Component nicht gefunden"})
		return
	}

	// Update basic fields
	component.Name = input.Name
	component.Description = input.Description
	component.Status = input.Status
	component.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := tx.Save(&component).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove existing tags
	if err := tx.Delete(&ComponentTag{}, "component_id = ?", component.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add new tags
	for _, tagName := range input.TagNames {
		var tag Tag

		// Try to find existing tag
		err := tx.Where("name = ?", tagName).First(&tag).Error
		if err != nil {
			// Create new tag if it doesn't exist
			tag = Tag{
				ID:   uuid.New().String(),
				Name: tagName,
			}
			if err := tx.Create(&tag).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Add tag to component
		if err := tx.Create(&ComponentTag{ComponentID: component.ID, TagID: tag.ID}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	// Fetch updated component with tags
	db.Preload("Tags").First(&component, "id = ?", id)

	c.JSON(http.StatusOK, component)
}

func deleteComponent(c *gin.Context) {
	id := c.Param("id")

	tx := db.Begin()

	// Delete component tags first
	if err := tx.Delete(&ComponentTag{}, "component_id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete component
	if err := tx.Delete(&Component{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.Status(http.StatusNoContent)
}
