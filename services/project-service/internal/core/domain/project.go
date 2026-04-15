package domain

import (
	"time"

	"github.com/google/uuid"
)

// Project — pure domain entity (no framework dependency, Kliops pattern).
type Project struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	IndustryType   string     `json:"industry_type"`
	Location       string     `json:"location"`
	Budget         float64    `json:"budget"`
	FloorWidth     float64    `json:"floor_width"`
	FloorDepth     float64    `json:"floor_depth"`
	TargetCapacity *string    `json:"target_capacity"`
	Status         string     `json:"status"`
	Version        int        `json:"version"`
	ArchivedAt     *time.Time `json:"archived_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewProject creates a Project with sensible defaults.
func NewProject(name, industryType, location string, budget, floorWidth, floorDepth float64, targetCapacity *string) Project {
	now := time.Now().UTC()
	return Project{
		ID:             uuid.New().String(),
		Name:           name,
		IndustryType:   industryType,
		Location:       location,
		Budget:         budget,
		FloorWidth:     floorWidth,
		FloorDepth:     floorDepth,
		TargetCapacity: targetCapacity,
		Status:         "active",
		Version:        1,
		ArchivedAt:     nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// IsArchived returns true when the project has been soft-deleted.
func (p *Project) IsArchived() bool {
	return p.ArchivedAt != nil
}

// SoftDelete marks the project as archived.
func (p *Project) SoftDelete() {
	now := time.Now().UTC()
	p.ArchivedAt = &now
	p.Status = "archived"
	p.UpdatedAt = now
}

// Restore un-archives the project.
func (p *Project) Restore() {
	p.ArchivedAt = nil
	p.Status = "active"
	p.UpdatedAt = time.Now().UTC()
}

// IncrementVersion bumps version and updated_at.
func (p *Project) IncrementVersion() {
	p.Version++
	p.UpdatedAt = time.Now().UTC()
}

// ApplyUpdate applies a partial update map, protecting immutable fields.
func (p *Project) ApplyUpdate(fields map[string]interface{}) {
	protected := map[string]bool{"id": true, "created_at": true, "status": true, "archived_at": true}

	for key, val := range fields {
		if protected[key] {
			continue
		}
		switch key {
		case "name":
			if v, ok := val.(string); ok {
				p.Name = v
			}
		case "industry_type":
			if v, ok := val.(string); ok {
				p.IndustryType = v
			}
		case "location":
			if v, ok := val.(string); ok {
				p.Location = v
			}
		case "budget":
			if v, ok := toFloat64(val); ok {
				p.Budget = v
			}
		case "floor_width":
			if v, ok := toFloat64(val); ok {
				p.FloorWidth = v
			}
		case "floor_depth":
			if v, ok := toFloat64(val); ok {
				p.FloorDepth = v
			}
		case "target_capacity":
			if val == nil {
				p.TargetCapacity = nil
			} else if v, ok := val.(string); ok {
				p.TargetCapacity = &v
			}
		}
	}
	p.IncrementVersion()
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}
