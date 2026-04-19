package main

import (
	"log"
	"net/http"

	"api-on/internal/auth"
	"api-on/internal/shared/config"
	"api-on/internal/shared/database"
	"api-on/internal/shared/middleware"
	"api-on/internal/shared/response"
	"api-on/internal/shared/security"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
	jwtpkg "api-on/pkg/jwt"
	"api-on/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.New()
	store := database.NewStore(cfg.DataFile)
	if err := store.Initialize(); err != nil {
		log.Fatal(err)
	}

	jwtSvc := jwtpkg.NewJWTService(cfg.SecretKey, cfg.JWTIssuer)

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := user.NewRepository(store)

	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtSvc)
	authHandler := auth.NewHandler(authUsecase)

	tenantHandler := tenant.NewHandler(tenant.NewUsecase(tenantRepo))
	subscriptionHandler := subscription.NewHandler(subscription.NewUsecase(subscriptionRepo))
	userHandler := user.NewHandler(user.NewUsecase(userRepo, subscriptionRepo))

	router := gin.Default()

	router.GET("/v1/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"message": "PsicoSistem backend is running",
			"env":     cfg.AppEnv,
		}, nil)
	})

	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", middleware.AuthRequired(jwtSvc, security.UserTypeInternal), authHandler.Refresh)
		}

		internalGroup := v1.Group("")
		internalGroup.Use(middleware.AuthRequired(jwtSvc, security.UserTypeInternal))
		{
			tenantGroup := internalGroup.Group("/tenant")
			{
				tenantGroup.GET("/me", tenantHandler.Me)
				tenantGroup.GET("/subscription", subscriptionHandler.Current)
			}

			userGroup := internalGroup.Group("/users")
			{
				userGroup.GET("/me", userHandler.Me)
				userGroup.GET("", userHandler.List)
				userGroup.POST("", userHandler.Create)
				userGroup.PATCH("/:id", userHandler.Update)
				userGroup.DELETE("/:id", userHandler.Delete)
			}
		}
	}

	appLogger.Info("starting API server", "port", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
