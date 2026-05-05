package main

import (
	"context"
	"log"
	"net/http"

	"api-on/internal/auth"
	"api-on/internal/learner"
	"api-on/internal/permission"
	"api-on/internal/shared/bootstrap"
	"api-on/internal/shared/config"
	"api-on/internal/shared/infra"
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
	runtimeProfile, err := infra.NewRuntimeProfile(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repositories, err := bootstrap.BuildRepositories(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer repositories.Close()

	jwtSvc := jwtpkg.NewJWTService(cfg.SecretKey, cfg.JWTIssuer)

	tenantRepo := repositories.TenantRepo
	subscriptionRepo := repositories.SubscriptionRepo
	userRepo := repositories.UserRepo
	learnerRepo := repositories.LearnerRepo

	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtSvc)
	authHandler := auth.NewHandler(authUsecase)

	tenantHandler := tenant.NewHandler(tenant.NewUsecase(tenantRepo))
	subscriptionHandler := subscription.NewHandler(subscription.NewUsecase(subscriptionRepo))
	userHandler := user.NewHandler(user.NewUsecase(userRepo, subscriptionRepo))
	learnerHandler := learner.NewHandler(learner.NewUsecase(learnerRepo, subscriptionRepo))
	permissionHandler := permission.NewHandler(permission.NewUsecase(userRepo))
	identityResolver := user.NewIdentityResolver(userRepo)

	router := gin.Default()

	router.GET("/v1/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"message":        "PsicoSistem backend is running",
			"env":            cfg.AppEnv,
			"storage_driver": cfg.StorageDriver,
			"cloud_provider": runtimeProfile.CloudProvider,
		}, nil)
	})

	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", middleware.AuthRequiredWithResolver(jwtSvc, identityResolver, security.UserTypeInternal), authHandler.Refresh)
		}

		internalGroup := v1.Group("")
		internalGroup.Use(middleware.AuthRequiredWithResolver(jwtSvc, identityResolver, security.UserTypeInternal))
		{
			tenantGroup := internalGroup.Group("/tenant")
			{
				tenantGroup.GET("/me", tenantHandler.Me)
				tenantGroup.GET("/subscription", subscriptionHandler.Current)
			}

			permissionGroup := internalGroup.Group("/permissions")
			{
				permissionGroup.GET("/me", permissionHandler.Me)
			}

			userGroup := internalGroup.Group("/users")
			{
				userGroup.GET("/me", userHandler.Me)
				userGroup.GET("", userHandler.List)
				userGroup.POST("", userHandler.Create)
				userGroup.GET("/:id/permissions", permissionHandler.GetByUser)
				userGroup.PATCH("/:id/permissions", permissionHandler.UpdateByUser)
				userGroup.PATCH("/:id", userHandler.Update)
				userGroup.DELETE("/:id", userHandler.Delete)
			}

			learnerGroup := internalGroup.Group("/learners")
			{
				learnerGroup.GET("", learnerHandler.List)
				learnerGroup.POST("", learnerHandler.Create)
				learnerGroup.GET("/:id", learnerHandler.Get)
				learnerGroup.PATCH("/:id", learnerHandler.Update)
				learnerGroup.DELETE("/:id", learnerHandler.Delete)
			}
		}
	}

	appLogger.Info(
		"starting API server",
		"port", cfg.Port,
		"storage_driver", cfg.StorageDriver,
		"cloud_provider", runtimeProfile.CloudProvider,
	)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
