package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rmscoal/moilerplate/internal/app/service"
	"github.com/rmscoal/moilerplate/internal/app/usecase"
	"github.com/rmscoal/moilerplate/internal/delivery/middleware"
)

type CredentialController struct {
	BaseControllerV1
	uc  usecase.ICredentialUseCase
	svc service.IRaterService
}

func NewCredentialController(rg *gin.RouterGroup, uc usecase.ICredentialUseCase, svc service.IRaterService) {
	controller := new(CredentialController)
	controller.uc = uc
	controller.svc = svc

	r := rg.Group("/credentials")
	{
		// Rate limiting middleware - For all credentials endpoint group
		r.Use(middleware.NewMiddleware().RateLimiterMiddleware(controller.svc))

		r.POST("/signup", controller.signupHandler)
		r.POST("/login", controller.loginHandler)
		r.POST("/refresh", controller.refreshHandler)
	}
}

/*
*************************************************
Controllers
*************************************************
*/
func (controller *CredentialController) signupHandler(c *gin.Context) {
}

func (controller *CredentialController) loginHandler(c *gin.Context) {
}

func (controller *CredentialController) refreshHandler(c *gin.Context) {
}
