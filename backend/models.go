package main

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:'operational'"` // operational, degraded, outage
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Incident struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // investigating, identified, monitoring, resolved
	Impact      string    `json:"impact"` // critical, major, minor
	ServiceID   uuid.UUID `json:"service_id"`
	Service     Service   `json:"service"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

type IncidentUpdate struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	IncidentID  uuid.UUID `json:"incident_id"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
