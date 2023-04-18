package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto/mapper"
)

type CredentialController struct {
	// V1BaseController
	uc usecase.ICredentialUseCase
}

func NewCredentialController(rg *gin.RouterGroup, uc usecase.ICredentialUseCase) {
	controller := new(CredentialController)
	controller.uc = uc

	r := rg.Group("/credentials")
	{
		r.POST("/signup", controller.signupHandler)
	}
}

/*
*************************************************
Controllers
*************************************************
*/
func (controller *CredentialController) signupHandler(c *gin.Context) {
	var raw dto.SignUpRequest
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		// controller.ClientError(c, err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	dto := mapper.SignUpRequestToUserDomain(raw)
	data, err := controller.uc.SignUp(c.Request.Context(), dto)
	if err != nil {
		panic(err)
	}
	c.JSON(201, data)
}
