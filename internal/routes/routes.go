package routes

import (
	"net/http"

	"common/pkg/config"
	"common/pkg/database/mongodb"
	"common/pkg/handlers"
	"makatom-api-config/internal/models"
	"makatom-api-config/internal/services"
)

// RegisterConfigRoutes sets up and returns the main router for the config service.
// It now uses the custom GenericRouter to handle dynamic path parameters.
func RegisterConfigRoutes() http.Handler {
	cfg := config.GetConfig()

	// 1. Use the new GenericRouter instead of the standard ServeMux.
	mux := http.NewServeMux()

	// Get MongoDB connection
	client, _ := mongodb.Manager.Get(cfg.MongoURIName)
	db := client.Database(cfg.MongoDatabase)
	collection := db.Collection("configs")

	// Create service directly
	configService := services.NewConfigService(collection)

	// Define APIs directly using service functions.
	// The Path field now includes the HTTP method, which the new router uses.
	apis := []handlers.APIDefinition{
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
	}

	// 2. Register the routes with the new GenericRouter instance.
	// No other changes are needed here.
	handlers.RegisterRoutes(mux, apis)

	return mux
}
