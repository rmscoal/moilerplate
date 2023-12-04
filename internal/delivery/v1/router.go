package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/composer"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/middleware"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/logger"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(r *gin.Engine, logger *logger.AppLogger, ucComposer composer.IUseCaseComposer) {
	r.Use(middleware.LogRequestMiddleware(logger))
	r.LoadHTMLGlob("public/**/*")

	// API V1 - Parent of all endpoint for V1.
	v1 := r.Group("/api/v1")

	// Credentials controller
	NewCredentialController(v1, ucComposer.CredentialUseCase(), ucComposer.RaterUseCase())
	// Admin controller
	NewAdminController(v1, ucComposer.CredentialUseCase())

	// Protected endpoint
	ptd := v1.Group("/ptd")
	// Authenticate middleware - For all protected endpoint
	ptd.Use(middleware.NewMiddleware().AuthMiddleware(ucComposer.CredentialUseCase()))
	// Profile controller
	NewUserProfileController(ptd, ucComposer.UserProfileUseCase())
}
