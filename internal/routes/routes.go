package routes

import (
	"net/http"

	"common/pkg/config"
	"common/pkg/database/mongodb"
	"common/pkg/handlers"
	"makatom-api-config/internal/models"
	"makatom-api-config/internal/services"
)

func RegisterConfigRoutes() *http.ServeMux {
	cfg := config.GetConfig()

	r := http.NewServeMux()

	// Get MongoDB connection
	client, _ := mongodb.Manager.Get(cfg.MongoURIName)
	db := client.Database(cfg.MongoDatabase)
	collection := db.Collection("configs")

	// Create service directly
	configService := services.NewConfigService(collection)

	// Define APIs directly using service functions
	apis := []handlers.APIDefinition{
		// Create config
		{Path: "POST /config", Handler: handlers.GenerateHandler(configService.CreateConfig, new(models.CreateConfigRequest))},

		// Get all configs
		{Path: "GET /configs", Handler: handlers.GenerateHandler(configService.GetConfigs, new(models.ConfigQuery))},

		// Get config by ID
		{Path: "GET /config/get", Handler: handlers.GenerateHandler(configService.GetConfigByID, new(models.ConfigIDRequest))},

		// Update config
		{Path: "PUT /config/update", Handler: handlers.GenerateHandler(configService.UpdateConfig, new(models.UpdateConfigWithIDRequest))},

		// Delete config
		{Path: "DELETE /config/delete", Handler: handlers.GenerateHandler(configService.DeleteConfig, new(models.ConfigIDRequest))},
	}

	handlers.RegisterRoutes(r, apis)

	return r
}
