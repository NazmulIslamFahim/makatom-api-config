package models

import (
	"time"

	"common/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Config represents a configuration entity
type Config struct {
	*types.Base
	Name          string                 `bson:"name" json:"name" validate:"required"`
	Type          string                 `bson:"type" json:"type" validate:"required"`
	Subtype       string                 `bson:"subtype" json:"subtype,omitempty"`
	Tags          []string               `bson:"tags" json:"tags,omitempty"`
	TenantID      string                 `bson:"tenant_id" json:"tenant_id" validate:"required"`
	CreatedBy     string                 `bson:"created_by" json:"created_by" validate:"required"`
	LastUpdatedBy string                 `bson:"last_updated_by" json:"last_updated_by" validate:"required"`
	Metadata      map[string]interface{} `bson:"metadata" json:"metadata,omitempty"`
}

// CreateConfigRequest represents the request payload for creating a config
type CreateConfigRequest struct {
	Name     string                 `json:"name" validate:"required"`
	Type     string                 `json:"type" validate:"required"`
	Subtype  string                 `json:"subtype,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	TenantID string                 `json:"tenant_id" validate:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateConfigRequest represents the request payload for updating a config
type UpdateConfigRequest struct {
	Name     string                 `json:"name,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Subtype  string                 `json:"subtype,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateConfigWithIDRequest represents the request payload for updating a config with ID
type UpdateConfigWithIDRequest struct {
	ID       string                 `query:"id" validate:"required"`
	Name     string                 `json:"name,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Subtype  string                 `json:"subtype,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ConfigQuery represents query parameters for filtering configs
type ConfigQuery struct {
	TenantID string `query:"tenant_id" validate:"required"`
	Type     string `query:"type,omitempty"`
	Subtype  string `query:"subtype,omitempty"`
	Tag      string `query:"tag,omitempty"`
	Limit    int64  `query:"limit,omitempty"`
	Skip     int64  `query:"skip,omitempty"`
}

// ConfigIDRequest represents request with config ID from path
type ConfigIDRequest struct {
	ID string `query:"id" validate:"required"`
}

// ConfigResponse represents the response payload for config operations
type ConfigResponse struct {
	ID            primitive.ObjectID     `json:"id"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	Subtype       string                 `json:"subtype,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	TenantID      string                 `json:"tenant_id"`
	CreatedBy     string                 `json:"created_by"`
	LastUpdatedBy string                 `json:"last_updated_by"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// ToResponse converts a Config to ConfigResponse
func (c *Config) ToResponse() ConfigResponse {
	return ConfigResponse{
		ID:            c.ID,
		Name:          c.Name,
		Type:          c.Type,
		Subtype:       c.Subtype,
		Tags:          c.Tags,
		TenantID:      c.TenantID,
		CreatedBy:     c.CreatedBy,
		LastUpdatedBy: c.LastUpdatedBy,
		Metadata:      c.Metadata,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}
