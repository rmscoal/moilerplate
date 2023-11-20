package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto/mapper"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/model"
)

type CredentialController struct {
	model.BaseControllerV1
	uc usecase.ICredentialUseCase
}

func NewCredentialController(rg *gin.RouterGroup, uc usecase.ICredentialUseCase) {
	controller := new(CredentialController)
	controller.uc = uc

	r := rg.Group("/credentials")
	{
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
	var raw dto.SignUpRequest
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		controller.ClientError(c, err)
		return
	}

	req := mapper.Credential.SignUpRequestToUserDomain(raw)
	user, err := controller.uc.SignUp(c.Request.Context(), req)
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Created(c, mapper.Credential.UserDomainToTokenResponse(user))
}

func (controller *CredentialController) loginHandler(c *gin.Context) {
	var raw dto.LoginRequest
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		controller.ClientError(c, err)
		return
	}

	req := mapper.Credential.LoginRequestToUserCredential(raw)
	user, err := controller.uc.Login(c.Request.Context(), req)
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Ok(c, mapper.Credential.UserDomainToTokenResponse(user))
}

func (controller *CredentialController) refreshHandler(c *gin.Context) {
	var raw dto.RefreshRequest
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		controller.ClientError(c, err)
		return
	}

	user, err := controller.uc.Refresh(c.Request.Context(), raw.RefreshToken)
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Created(c, mapper.Credential.UserDomainToTokenResponse(user))
}
