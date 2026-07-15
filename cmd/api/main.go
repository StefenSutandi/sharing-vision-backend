package main

import (
	"log"
	"net/http"

	"github.com/StefenSutandi/sharing-vision-backend/internal/config"
	"github.com/StefenSutandi/sharing-vision-backend/internal/handler"
	"github.com/StefenSutandi/sharing-vision-backend/internal/repository"
	"github.com/StefenSutandi/sharing-vision-backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	r := gin.Default()

	configCORS := cors.DefaultConfig()
	configCORS.AllowAllOrigins = true
	configCORS.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	configCORS.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(configCORS))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	articleGroup := r.Group("/article")
	{
		articleGroup.POST("/", articleHandler.Create)
		articleGroup.POST("", articleHandler.Create)
		articleGroup.GET("/:limit/:offset", articleHandler.List)
		articleGroup.GET("/:id", articleHandler.GetByID)
		articleGroup.PUT("/:id", articleHandler.Update)
		articleGroup.DELETE("/:id", articleHandler.Delete)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
