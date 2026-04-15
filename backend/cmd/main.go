package main

import (
	"os"

	"api-on/internal/auth/handler"
	"api-on/internal/organization/middleware"
	"api-on/internal/organization/repository"
	"api-on/internal/auth/usecase"
	jwtpkg "api-on/pkg/jwt"


	"github.com/gin-gonic/gin"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-key-change-this"
	}

	userRepo := repository.NewMemoryUserRepository()
	jwtSvc := jwtpkg.NewJWTService(secret, "go-backend")

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtSvc)
	userUsecase := usecase.NewUserUsecase(userRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "API is running"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/users")
		protected.Use(middleware.AuthRequired(jwtSvc))
		{
			protected.GET("/me", userHandler.Me)
			protected.GET("", userHandler.List)
			protected.POST("", userHandler.Create)
			protected.PUT("/:id", userHandler.Update)
			protected.DELETE("/:id", userHandler.Delete)
		}
	}

	r.Run(":8080")
}