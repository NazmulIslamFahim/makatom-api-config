package services

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"makatom-api-config/internal/models"
	"makatom/common/pkg/database/mongodb"
	"makatom/common/pkg/handlers"
	"makatom/common/pkg/types"
)

const (
	MaxArchiveHistory = 10
)

// ConfigService handles business logic for config operations
type ConfigService struct {
	repo        *mongodb.MongoRepository[models.Config]
	archiveRepo *mongodb.MongoRepository[models.ConfigArchive]
}

// NewConfigService creates a new ConfigService instance
func NewConfigService(configCollection, archiveCollection *mongo.Collection) *ConfigService {
	return &ConfigService{
		repo:        mongodb.NewMongoRepository[models.Config](configCollection),
		archiveRepo: mongodb.NewMongoRepository[models.ConfigArchive](archiveCollection),
	}
}

// CreateConfig creates a new config
func (s *ConfigService) CreateConfig(ctx context.Context, req models.CreateConfigRequest) handlers.ServiceResponse {
	// For now, use dummy values - in a real app, these would come from JWT token
	userID := "dummy-user-id"
	tenantID := "dummy-tenant-id"

	// Validate that type exists
	_, typeExists := types.GlobalConfigTypeRegistry.GetType(req.Type)
	if !typeExists {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "config type does not exist",
		}
	}

	// Validate that subtype exists for the type
	if req.Subtype != "" {
		_, subtypeExists := types.GlobalConfigTypeRegistry.GetSubtype(req.Type, req.Subtype)
		if !subtypeExists {
			return handlers.ServiceResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "config subtype does not exist for the given type",
			}
		}
	}

	// Validate metadata against subtype schema if metadata is provided
	if req.Metadata != nil {
		validationResult := types.GlobalConfigTypeRegistry.ValidateMetadata(req.Type, req.Subtype, req.Metadata)
		if !validationResult.Valid {
			return handlers.ServiceResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "metadata validation failed",
				Data:       validationResult,
			}
		}
	}

	// Check if config with same name already exists for this tenant
	existing, err := s.repo.FindOne(ctx, bson.M{
		"name":      req.Name,
		"tenant_id": tenantID,
		"type":      req.Type,
		"subtype":   req.Subtype,
	})

	// If we found an existing config, return duplicate error
	if err == nil && existing.ID != primitive.NilObjectID {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "config with this name already exists for this tenant and type",
		}
	}

	// If we got a "not found" error, that's good - proceed
	// If we got any other error, return database error
	if err != nil && err.Error() != "not found" {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to check existing config: %v", err),
		}
	}

	// Create new config
	config := models.Config{
		Base:          &types.Base{},
		Name:          req.Name,
		Type:          req.Type,
		Subtype:       req.Subtype,
		Tags:          req.Tags,
		TenantID:      tenantID,
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
	// For now, use dummy tenant ID - in a real app, this would come from JWT token
	tenantID := "dummy-tenant-id"

	// Build filter
	filter := bson.M{"tenant_id": tenantID}

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

// UpdateConfig updates an existing config with transaction support
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

	// Do not allow changing name, type, subtype, or tenantID (before DB lookup)
	if req.Name != "" {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "changing config name is not allowed",
		}
	}
	if req.Type != "" {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "changing config type is not allowed",
		}
	}
	if req.Subtype != "" {
		return handlers.ServiceResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "changing config subtype is not allowed",
		}
	}
	// tenantID is not updatable by design (not in request struct)

	// Get the existing config by id and tenantID
	existing, err := s.repo.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID})
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

	// Validate metadata against subtype schema if metadata is being updated
	if req.Metadata != nil {
		validationResult := types.GlobalConfigTypeRegistry.ValidateMetadata(existing.Type, existing.Subtype, req.Metadata)
		if !validationResult.Valid {
			return handlers.ServiceResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "metadata validation failed",
				Data:       validationResult,
			}
		}
	}

	// Build update document (only allow tags and metadata)
	updates := bson.M{}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}
	updates["last_updated_by"] = userID

	var updatedConfig models.Config

	// Use transaction to ensure both archive creation and config update happen atomically
	err = s.repo.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
		// Archive the current version before updating
		err := s.archiveConfigVersionWithSession(sessCtx, existing, userID)
		if err != nil {
			return fmt.Errorf("failed to archive config version: %v", err)
		}

		// Update the config
		updated, err := s.repo.UpdateByID(sessCtx, id, bson.M{"$set": updates})
		if err != nil {
			return fmt.Errorf("failed to update config: %v", err)
		}
		updatedConfig = updated
		return nil
	})

	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("transaction failed: %v", err),
		}
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusOK,
		Data:       updatedConfig.ToResponse(),
	}
}

// DeleteConfig deletes a config by its ID and all its archives with transaction support
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

	// Check if config exists and belongs to tenant
	existing, err := s.repo.FindByID(ctx, id)
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

	if existing.TenantID != tenantID {
		return handlers.ServiceResponse{
			StatusCode: http.StatusNotFound,
			Error:      "Config not found",
		}
	}

	// Use transaction to ensure both archive deletion and config deletion happen atomically
	err = s.repo.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
		// Delete all archives for this config first
		err := s.deleteAllArchivesByConfigIDWithSession(sessCtx, id)
		if err != nil {
			return fmt.Errorf("failed to delete config archives: %v", err)
		}

		// Delete the config
		_, err = s.repo.FindOneAndDelete(sessCtx, bson.M{
			"_id":       id,
			"tenant_id": tenantID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete config: %v", err)
		}
		return nil
	})

	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("transaction failed: %v", err),
		}
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusOK,
		Data:       map[string]string{"message": "Config and all archives deleted successfully"},
	}
}

// GetConfigArchives retrieves archive history for a config
func (s *ConfigService) GetConfigArchives(ctx context.Context, req models.ConfigIDRequest) handlers.ServiceResponse {
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

	// Check if config exists and belongs to tenant
	existing, err := s.repo.FindByID(ctx, id)
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

	if existing.TenantID != tenantID {
		return handlers.ServiceResponse{
			StatusCode: http.StatusNotFound,
			Error:      "Config not found",
		}
	}

	// Get archives for this config, ordered by version descending
	archives, err := s.archiveRepo.Find(ctx, bson.M{
		"config_id": id,
		"tenant_id": tenantID,
	}, 0, 0) // No pagination for archives
	if err != nil {
		return handlers.ServiceResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("failed to get config archives: %v", err),
		}
	}

	// Convert to responses
	responses := make([]models.ConfigArchiveResponse, len(archives))
	for i, archive := range archives {
		responses[i] = archive.ToArchiveResponse()
	}

	return handlers.ServiceResponse{
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"archives": responses,
			"total":    len(responses),
		},
	}
}

// archiveConfigVersion archives the current version of a config
// func (s *ConfigService) archiveConfigVersion(ctx context.Context, config models.Config, archivedBy string) error {
// 	// Get current version number
// 	currentVersion, err := s.archiveRepo.Count(ctx, bson.M{"config_id": config.ID})
// 	if err != nil {
// 		return err
// 	}

// 	// Create archive entry
// 	archive := config.ToArchive(int(currentVersion)+1, archivedBy)
// 	_, err = s.archiveRepo.InsertOne(ctx, archive)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if we need to remove oldest archive (keep only MaxArchiveHistory)
// 	// Since we just added one archive, if currentVersion + 1 > MaxArchiveHistory, we need to remove the oldest
// 	if currentVersion+1 > MaxArchiveHistory {
// 		// Find the oldest archive (lowest version)
// 		oldestArchive, err := s.archiveRepo.FindOne(ctx, bson.M{
// 			"config_id": config.ID,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		// Delete the oldest archive
// 		_, err = s.archiveRepo.DeleteOne(ctx, bson.M{"_id": oldestArchive.ID})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// deleteAllArchivesByConfigID deletes all archives for a specific config ID
// func (s *ConfigService) deleteAllArchivesByConfigID(ctx context.Context, configID primitive.ObjectID) error {
// 	// Delete all archives for this config
// 	_, err := s.archiveRepo.DeleteMany(ctx, bson.M{"config_id": configID})
// 	return err
// }

// archiveConfigVersionWithSession archives the current version of a config within a session
func (s *ConfigService) archiveConfigVersionWithSession(sessCtx mongo.SessionContext, config models.Config, archivedBy string) error {
	// Get current version number
	currentVersion, err := s.archiveRepo.Count(sessCtx, bson.M{"config_id": config.ID})
	if err != nil {
		return err
	}

	// Create archive entry
	archive := config.ToArchive(int(currentVersion)+1, archivedBy)
	_, err = s.archiveRepo.InsertOne(sessCtx, archive)
	if err != nil {
		return err
	}

	// Check if we need to remove oldest archive (keep only MaxArchiveHistory)
	// Since we just added one archive, if currentVersion + 1 > MaxArchiveHistory, we need to remove the oldest
	if currentVersion+1 > MaxArchiveHistory {
		// Find the oldest archive (lowest version)
		oldestArchive, err := s.archiveRepo.FindOne(sessCtx, bson.M{
			"config_id": config.ID,
		})
		if err != nil {
			return err
		}

		// Delete the oldest archive
		_, err = s.archiveRepo.DeleteOne(sessCtx, bson.M{"_id": oldestArchive.ID})
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteAllArchivesByConfigIDWithSession deletes all archives for a specific config ID within a session
func (s *ConfigService) deleteAllArchivesByConfigIDWithSession(sessCtx mongo.SessionContext, configID primitive.ObjectID) error {
	// Delete all archives for this config
	_, err := s.archiveRepo.DeleteMany(sessCtx, bson.M{"config_id": configID})
	return err
}
