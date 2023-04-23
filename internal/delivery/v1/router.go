package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/composer"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/middleware"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/logger"
)

func NewRouter(r *gin.Engine, logger *logger.AppLogger, ucComposer composer.IUseCaseComposer) {
	r.Use(middleware.LogRequestMiddleware(logger))

	v1 := r.Group("/api/v1")
	{
		// Rate limiting middleware
		v1.Use(middleware.NewMiddleware().RateLimiterMiddleware(ucComposer.RaterUseCase()))

		NewCredentialController(v1, ucComposer.CredentialUseCase())

		ptd := v1.Group("/ptd")
		{
			// Authorizations middleware
			ptd.Use(middleware.NewMiddleware().AuthMiddleware(ucComposer.CredentialUseCase()))

			NewUserProfileController(ptd, ucComposer.UserProfileUseCase())
		}
	}
}
