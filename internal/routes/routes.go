package routes

import (
	"net/http"

	"makatom-api-config/internal/models"
	configServices "makatom-api-config/internal/services"
	"makatom/common/pkg/config"
	"makatom/common/pkg/database/mongodb"
	"makatom/common/pkg/handlers"
	commonServices "makatom/common/pkg/services"
	"makatom/common/pkg/types"
)

// RegisterConfigRoutes sets up and returns the main router for the config service.
// It now uses the custom GenericRouter to handle dynamic path parameters.
func RegisterConfigRoutes() http.Handler {
	cfg := config.GetConfig()

	// Initialize the type system
	types.Init()

	// 1. Use the new GenericRouter instead of the standard ServeMux.
	mux := http.NewServeMux()

	// Get MongoDB connection
	client, _ := mongodb.Manager.Get(cfg.MongoURIName)
	db := client.Database(cfg.MongoDatabase)
	configCollection := db.Collection("configs")
	archiveCollection := db.Collection("config_archives")

	// Create services
	configService := configServices.NewConfigService(configCollection, archiveCollection)
	configTypeService := commonServices.NewConfigTypeService()

	// Define APIs directly using service functions.
	// The Path field now includes the HTTP method, which the new router uses.
	apis := []handlers.APIDefinition{
		// Config APIs
		// Create config
		{
			Path:    "POST /config",
			Handler: handlers.GenerateHandler(configService.CreateConfig, new(models.CreateConfigRequest)),
		},

		// Get all configs
		{
			Path:    "GET /configs",
			Handler: handlers.GenerateHandler(configService.GetConfigs, new(models.ConfigQuery)),
		},

		// Get config by ID
		{
			Path:    "GET /config",
			Handler: handlers.GenerateHandler(configService.GetConfigByID, new(models.ConfigIDRequest)),
		},

		// Update config
		{
			Path:    "PUT /config",
			Handler: handlers.GenerateHandler(configService.UpdateConfig, new(models.UpdateConfigWithIDRequest)),
		},

		// Delete config
		{
			Path:    "DELETE /config",
			Handler: handlers.GenerateHandler(configService.DeleteConfig, new(models.ConfigIDRequest)),
		},

		// Get config archives
		{
			Path:    "GET /config/archives",
			Handler: handlers.GenerateHandler(configService.GetConfigArchives, new(models.ConfigIDRequest)),
		},

		// Type APIs
		// Get all types
		{
			Path:    "GET /types",
			Handler: handlers.GenerateHandler[types.EmptyRequest](configTypeService.GetAllTypes, new(types.EmptyRequest)),
		},

		// Get specific type
		{
			Path:    "GET /types/{type}",
			Handler: handlers.GenerateHandler(configTypeService.GetType, new(types.TypeRequest)),
		},

		// Get subtypes for a type
		{
			Path:    "GET /types/{type}/subtypes",
			Handler: handlers.GenerateHandler(configTypeService.GetSubtypes, new(types.TypeRequest)),
		},

		// Get specific subtype
		{
			Path:    "GET /types/{type}/subtypes/{subtype}",
			Handler: handlers.GenerateHandler(configTypeService.GetSubtype, new(types.SubtypeRequest)),
		},

		// Validate metadata
		{
			Path:    "POST /validate-metadata",
			Handler: handlers.GenerateHandler(configTypeService.ValidateMetadata, new(types.ValidationRequest)),
		},

		// Decrypt config field
		{
			Path:    "POST /config/decrypt",
			Handler: handlers.GenerateHandler(configService.DecryptConfigField, new(models.DecryptFieldRequest)),
		},
	}

	// 2. Register the routes with the new GenericRouter instance.
	// No other changes are needed here.
	handlers.RegisterRoutes(mux, apis)

	return mux
}
