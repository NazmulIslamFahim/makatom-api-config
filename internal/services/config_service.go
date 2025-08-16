package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"common/pkg/database/mongodb"
	"common/pkg/handlers"
	"common/pkg/types"
	"makatom-api-config/internal/models"
)

// ConfigService handles business logic for config operations
type ConfigService struct {
	repo *mongodb.MongoRepository[models.Config]
}

// NewConfigService creates a new ConfigService instance
func NewConfigService(collection *mongo.Collection) *ConfigService {
	return &ConfigService{
		repo: mongodb.NewMongoRepository[models.Config](collection),
	}
}

// CreateConfig creates a new config
func (s *ConfigService) CreateConfig(ctx context.Context, req models.CreateConfigRequest) handlers.ServiceResponse {
	// For now, use dummy user ID - in a real app, this would come from JWT token
	userID := "dummy-user-id"

	// Add tenant_id to the request if not provided (for demo purposes)
	if req.TenantID == "" {
		req.TenantID = "dummy-tenant-id"
	}

	// Check if config with same name already exists for this tenant
	existing, err := s.repo.FindOne(ctx, bson.M{
		"name":      req.Name,
		"tenant_id": req.TenantID,
		"type":      req.Type,
	})

	if err != nil && !errors.Is(err, errors.New("not found")) {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to check existing config: %v", err),
		}
	}

	if existing.ID != primitive.NilObjectID {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "config with this name already exists for this tenant and type",
		}
	}

	// Create new config
	config := models.Config{
		Base:          &types.Base{},
		Name:          req.Name,
		Type:          req.Type,
		Subtype:       req.Subtype,
		Tags:          req.Tags,
		TenantID:      req.TenantID,
		CreatedBy:     userID,
		LastUpdatedBy: userID,
		Metadata:      req.Metadata,
	}

	createdConfig, err := s.repo.InsertOne(ctx, config)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("failed to create config: %v", err),
		}
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusCreated,
		Data:       createdConfig.ToResponse(),
	}
}

// GetConfigByID retrieves a config by its ID
func (s *ConfigService) GetConfigByID(ctx context.Context, req models.ConfigIDRequest) handlers.ServiceResponse {
	// Parse ObjectID
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid config ID",
		}
	}

	// For now, use dummy tenant ID - in a real app, this would come from JWT token
	tenantID := "dummy-tenant-id"

	config, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err.Error() == "not found" {
			return handlers.ServiceResponse{
				StatusCode: http.StatusNotFound,
				Error:      "Config not found",
			}
		}
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to get config: %v", err),
		}
	}

	// Ensure the config belongs to the requesting tenant
	if config.TenantID != tenantID {
		return handlers.ServiceResponse{
			StatusCode: http.StatusNotFound,
			Error:      "Config not found",
		}
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusOK,
		Data:       config.ToResponse(),
	}
}

// GetConfigs retrieves configs with filtering and pagination
func (s *ConfigService) GetConfigs(ctx context.Context, query models.ConfigQuery) handlers.ServiceResponse {
	// For now, use dummy tenant ID if not provided - in a real app, this would come from JWT token
	if query.TenantID == "" {
		query.TenantID = "dummy-tenant-id"
	}

	// Build filter
	filter := bson.M{"tenant_id": query.TenantID}

	if query.Type != "" {
		filter["type"] = query.Type
	}

	if query.Subtype != "" {
		filter["subtype"] = query.Subtype
	}

	if query.Tag != "" {
		filter["tags"] = bson.M{"$in": []string{query.Tag}}
	}

	// Get total count
	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to count configs: %v", err),
		}
	}

	// Set default limit if not provided
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get configs
	configs, err := s.repo.Find(ctx, filter, query.Skip, limit)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to get configs: %v", err),
		}
	}

	// Convert to responses
	responses := make([]models.ConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = config.ToResponse()
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"configs": responses,
			"total":   total,
			"limit":   query.Limit,
			"skip":    query.Skip,
		},
	}
}

// UpdateConfig updates an existing config
func (s *ConfigService) UpdateConfig(ctx context.Context, req models.UpdateConfigWithIDRequest) handlers.ServiceResponse {
	// Parse ObjectID
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid config ID",
		}
	}

	// For now, use dummy values - in a real app, these would come from JWT token
	tenantID := "dummy-tenant-id"
	userID := "dummy-user-id"

	// First, get the existing config to ensure it exists and belongs to the tenant
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to get config: %v", err),
		}
	}

	if existing.TenantID != tenantID {
		return handlers.ServiceResponse{
			StatusCode: http.StatusNotFound,
			Error:      "Config not found",
		}
	}

	// Check if name is being updated and if it conflicts with existing configs
	if req.Name != "" && req.Name != existing.Name {
		conflict, err := s.repo.FindOne(ctx, bson.M{
			"name":      req.Name,
			"tenant_id": tenantID,
			"type":      existing.Type,
			"_id":       bson.M{"$ne": id},
		})

		if err != nil && !errors.Is(err, errors.New("not found")) {
			return handlers.ServiceResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      fmt.Sprintf("failed to check name conflict: %v", err),
			}
		}

		if conflict.ID != primitive.NilObjectID {
			return handlers.ServiceResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "config with this name already exists for this tenant and type",
			}
		}
	}

	// Build update document
	updates := bson.M{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Subtype != "" {
		updates["subtype"] = req.Subtype
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}
	updates["last_updated_by"] = userID

	// Update the config
	updatedConfig, err := s.repo.UpdateByID(ctx, id, bson.M{"$set": updates})
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("failed to update config: %v", err),
		}
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusOK,
		Data:       updatedConfig.ToResponse(),
	}
}

// DeleteConfig deletes a config by its ID
func (s *ConfigService) DeleteConfig(ctx context.Context, req models.ConfigIDRequest) handlers.ServiceResponse {
	// Parse ObjectID
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid config ID",
		}
	}

	// For now, use dummy tenant ID - in a real app, this would come from JWT token
	tenantID := "dummy-tenant-id"

	// First, get the existing config to ensure it exists and belongs to the tenant
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to get config: %v", err),
		}
	}

	if existing.TenantID != tenantID {
		return handlers.ServiceResponse{
			StatusCode: http.StatusNotFound,
			Error:      "Config not found",
		}
	}

	// Delete the config
	_, err = s.repo.DeleteByID(ctx, id)
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to delete config: %v", err),
		}
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusNoContent,
		Data:       nil,
	}
}
