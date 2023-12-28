package v1

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/moilerplate/internal/composer"
	"github.com/rmscoal/moilerplate/internal/delivery/middleware"
	"github.com/rmscoal/moilerplate/pkg/logger"
)

//go:embed web/*
var web embed.FS

func NewRouter(r *gin.Engine, logger *logger.AppLogger, ucComposer composer.IUseCaseComposer) {
	// TODO: r.UseH2C = true
	r.Use(middleware.NewMiddleware().MetricsMiddleware())
	r.Use(middleware.NewMiddleware().TraceMiddleware())

	// Load all web html templates
	htmls := template.Must(template.ParseFS(web, "web/**/*.html"))

	r.Use(middleware.LogRequestMiddleware(logger))
	r.SetHTMLTemplate(htmls)

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
