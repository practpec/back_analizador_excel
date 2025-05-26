package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"analizador-backend/internal/application/services"
	"analizador-backend/internal/infrastructure/handlers"
	"analizador-backend/internal/infrastructure/repositories"
)

func main() {
	// Inicializar dependencias
	contactRepo := repositories.NewInMemoryContactRepository()
	validatorService := services.NewValidatorService()
	contactService := services.NewContactService(contactRepo, validatorService)
	contactHandler := handlers.NewContactHandler(contactService)

	// Configurar router
	router := gin.Default()
	
	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Rutas
	api := router.Group("/api/v1")
	{
		api.POST("/contacts/upload", contactHandler.UploadExcel)
		api.GET("/contacts", contactHandler.GetContacts)
		api.GET("/contacts/search", contactHandler.SearchContacts)
		api.PUT("/contacts/:id", contactHandler.UpdateContact)
		api.GET("/contacts/validate", contactHandler.ValidateContacts)
		api.GET("/contacts/download", contactHandler.DownloadExcel)
	}

	log.Println("Servidor iniciado en puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}