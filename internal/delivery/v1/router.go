package v1

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/moilerplate/internal/composer"
	"github.com/rmscoal/moilerplate/internal/delivery/middleware"
)

//go:embed web/*
var web embed.FS

func NewRouter(r *gin.Engine, ucComposer composer.IUseCaseComposer, svcComposer composer.IServiceComposer) {
	r.Use(middleware.NewMiddleware().MetricsMiddleware())
	r.Use(middleware.NewMiddleware().TraceMiddleware())

	// Load all web html templates
	htmls := template.Must(template.ParseFS(web, "web/**/*.html"))
	r.SetHTMLTemplate(htmls)

	// API V1 - Parent of all endpoint for V1.
	v1 := r.Group("/api/v1")

	// Credentials controller
	NewCredentialController(v1, ucComposer.CredentialUseCase(), svcComposer.RaterService())
}
