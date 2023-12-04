package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/middleware"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/model"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type AdminController struct {
	model.BaseControllerV1
	uc usecase.ICredentialUseCase
}

func NewAdminController(rg *gin.RouterGroup, uc usecase.ICredentialUseCase) {
	controller := new(AdminController)
	controller.uc = uc

	r := rg.Group("/docs")
	{
		r.GET("auth", controller.loginHandler)
		r.GET(":regex", middleware.NewMiddleware().AdminMiddleware(uc),
			ginSwagger.WrapHandler(
				swaggerFiles.Handler,
				ginSwagger.DefaultModelsExpandDepth(-1),
			),
		)
	}
}

func (controller *AdminController) loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", gin.H{})
}
