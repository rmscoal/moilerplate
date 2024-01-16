package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/moilerplate/internal/app/usecase"
	"github.com/rmscoal/moilerplate/internal/delivery/middleware"
	"github.com/rmscoal/moilerplate/internal/delivery/v1/model"

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
		r.GET("login", controller.loginPageHandler)
		r.POST("login", controller.loginHandler)
		r.GET(":regex", middleware.NewMiddleware().AdminMiddleware(controller.uc),
			ginSwagger.WrapHandler(
				swaggerFiles.Handler,
				ginSwagger.DefaultModelsExpandDepth(-1),
			),
		)
	}
}

func (controller *AdminController) loginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", gin.H{})
}

func (controller *AdminController) loginHandler(c *gin.Context) {
	adminKey := c.PostForm("key")

	session, err := controller.uc.AdminLogin(c.Request.Context(), adminKey)
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	c.SetCookie("x-session-key", session.Session, 3600, "/api/v1/docs", "localhost", true, true)
	c.Redirect(http.StatusFound, "index.html")
}
