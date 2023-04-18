package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/composer"
)

func NewRouter(r *gin.Engine, ucComposer composer.IUseCaseComposer) {
	r.Use()

	v1 := r.Group("/api/v1")
	{
		NewCredentialController(v1, ucComposer.CredentialUseCase())
	}
}
