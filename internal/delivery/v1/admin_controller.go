package v1

import (
	"net/http"
	"time"

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
		r.GET("login", controller.loginHandler)
		r.POST("verify", controller.verifyHandler)
		r.GET(":regex", middleware.NewMiddleware().AdminMiddleware(controller.uc),
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

func (controller *AdminController) verifyHandler(c *gin.Context) {
	adminKey := c.PostForm("key")

	if err := controller.uc.VerifyAdmin(c.Request.Context(), adminKey); err != nil {
		controller.Unauthorized(c, err)
		return
	}

	c.SetCookie("x-session-key", "ur mom", int(time.Now().Add(1*time.Hour).Unix()), "/api/v1/docs", "localhost", false, true)
	c.Redirect(http.StatusFound, "index.html")
}
