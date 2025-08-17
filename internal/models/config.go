package models

import (
	"time"

	"makatom/common/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Config represents a configuration entity
type Config struct {
	*types.Base   `bson:",inline"`
	Name          string                 `bson:"name" json:"name" validate:"required"`
	Type          string                 `bson:"type" json:"type" validate:"required"`
	Subtype       string                 `bson:"subtype" json:"subtype,omitempty"`
	Tags          []string               `bson:"tags" json:"tags,omitempty"`
	TenantID      string                 `bson:"tenant_id" json:"tenant_id" validate:"required"`
	CreatedBy     string                 `bson:"created_by" json:"created_by" validate:"required"`
	LastUpdatedBy string                 `bson:"last_updated_by" json:"last_updated_by" validate:"required"`
	Metadata      map[string]interface{} `bson:"metadata" json:"metadata,omitempty"`
}

// ConfigArchive represents a configuration archive entry
type ConfigArchive struct {
	*types.Base   `bson:",inline"`
	ConfigID      primitive.ObjectID     `bson:"config_id" json:"config_id"`
	Name          string                 `bson:"name" json:"name"`
	Type          string                 `bson:"type" json:"type"`
	Subtype       string                 `bson:"subtype" json:"subtype,omitempty"`
	Tags          []string               `bson:"tags" json:"tags,omitempty"`
	TenantID      string                 `bson:"tenant_id" json:"tenant_id"`
	CreatedBy     string                 `bson:"created_by" json:"created_by"`
	LastUpdatedBy string                 `bson:"last_updated_by" json:"last_updated_by"`
	Metadata      map[string]interface{} `bson:"metadata" json:"metadata,omitempty"`
	Version       int                    `bson:"version" json:"version"`
	ArchivedAt    time.Time              `bson:"archived_at" json:"archived_at"`
	ArchivedBy    string                 `bson:"archived_by" json:"archived_by"`
}

// CreateConfigRequest represents the request payload for creating a config
type CreateConfigRequest struct {
	Name     string                 `json:"name" validate:"required"`
	Type     string                 `json:"type" validate:"required"`
	Subtype  string                 `json:"subtype,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
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
	ID       string                 `param:"id" validate:"required"`
	Name     string                 `json:"name,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Subtype  string                 `json:"subtype,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ConfigQuery represents query parameters for filtering configs
type ConfigQuery struct {
	Type    string `param:"type,omitempty"`
	Subtype string `param:"subtype,omitempty"`
	Tag     string `param:"tag,omitempty"`
	Limit   int64  `param:"limit,omitempty"`
	Skip    int64  `param:"skip,omitempty"`
}

// ConfigIDRequest represents request with config ID from path
type ConfigIDRequest struct {
	ID string `param:"id" validate:"required"`
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

// ConfigArchiveResponse represents the response payload for config archive operations
type ConfigArchiveResponse struct {
	ID            primitive.ObjectID     `json:"id"`
	ConfigID      primitive.ObjectID     `json:"config_id"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	Subtype       string                 `json:"subtype,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	TenantID      string                 `json:"tenant_id"`
	CreatedBy     string                 `json:"created_by"`
	LastUpdatedBy string                 `json:"last_updated_by"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Version       int                    `json:"version"`
	ArchivedAt    time.Time              `json:"archived_at"`
	ArchivedBy    string                 `json:"archived_by"`
	CreatedAt     time.Time              `json:"created_at"`
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

// ToArchiveResponse converts a ConfigArchive to ConfigArchiveResponse
func (ca *ConfigArchive) ToArchiveResponse() ConfigArchiveResponse {
	return ConfigArchiveResponse{
		ID:            ca.ID,
		ConfigID:      ca.ConfigID,
		Name:          ca.Name,
		Type:          ca.Type,
		Subtype:       ca.Subtype,
		Tags:          ca.Tags,
		TenantID:      ca.TenantID,
		CreatedBy:     ca.CreatedBy,
		LastUpdatedBy: ca.LastUpdatedBy,
		Metadata:      ca.Metadata,
		Version:       ca.Version,
		ArchivedAt:    ca.ArchivedAt,
		ArchivedBy:    ca.ArchivedBy,
		CreatedAt:     ca.CreatedAt,
	}
}

// ToArchive converts a Config to ConfigArchive
func (c *Config) ToArchive(version int, archivedBy string) ConfigArchive {
	return ConfigArchive{
		Base:          &types.Base{},
		ConfigID:      c.ID,
		Name:          c.Name,
		Type:          c.Type,
		Subtype:       c.Subtype,
		Tags:          c.Tags,
		TenantID:      c.TenantID,
		CreatedBy:     c.CreatedBy,
		LastUpdatedBy: c.LastUpdatedBy,
		Metadata:      c.Metadata,
		Version:       version,
		ArchivedAt:    time.Now(),
		ArchivedBy:    archivedBy,
	}
}

// DecryptFieldRequest represents a request to decrypt a specific field
type DecryptFieldRequest struct {
	ConfigID  string `json:"config_id" validate:"required"`
	FieldName string `json:"field_name" validate:"required"`
}
