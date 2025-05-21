package main

type Tag struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
}

type Component struct {
	ID          string `json:"id" gorm:"type:uuid;primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'operational'"`
	Tags        []Tag  `json:"tags" gorm:"many2many:component_tags;"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ComponentTag struct {
	ComponentID string `gorm:"primaryKey"`
	TagID       string `gorm:"primaryKey"`
}
