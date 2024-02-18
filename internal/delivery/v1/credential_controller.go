package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rmscoal/moilerplate/internal/app/service"
	"github.com/rmscoal/moilerplate/internal/app/usecase"
	"github.com/rmscoal/moilerplate/internal/delivery/middleware"
	"github.com/rmscoal/moilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/moilerplate/internal/delivery/v1/dto/mapper"
	"github.com/rmscoal/moilerplate/internal/delivery/v1/model"
)

type CredentialController struct {
	model.BaseControllerV1
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

// SignupHandler godoc
//
//	@Summary		Sign up handler
//	@Description	Handles sign up for new users
//	@Tags			Credentials
//	@Accept			json
//	@Produce		json
//	@Param			signUpRequest	body		dto.SignUpRequest					true	"Signup request body"
//	@Success		200				{object}	model.Data{data=dto.SignUpResponse}
//	@Failure		409				{object}	model.Error{error=usecase.AppError}
//	@Failure		422				{object}	model.Error{error=usecase.AppError}
//	@Failure		500				{object}	model.Error{error=usecase.AppError}
//	@Router			/credentials/signup [post]
func (controller *CredentialController) signupHandler(c *gin.Context) {
	var req dto.SignUpRequest

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		controller.ClientError(c, usecase.NewClientError("Body", err))
		return
	}

	user, err := controller.uc.SignUp(c.Request.Context(), mapper.Credential.SignupRequestToUserDomain(req))
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Ok(c, mapper.Credential.UserDomainToSignUpResponse(user))
}

func (controller *CredentialController) loginHandler(c *gin.Context) {
}

func (controller *CredentialController) refreshHandler(c *gin.Context) {
}
