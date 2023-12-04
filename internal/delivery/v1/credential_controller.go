package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/service"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/middleware"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto/mapper"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/model"
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
//	@Summary		Sign up handler
//	@Description	Handles sign up for new users
//	@Tags			Credentials
//	@Accept			json
//	@Produce		json
//	@Param			signUpRequest	body		dto.SignUpRequest					true	"Signup request body"
//	@Success		200				{object}	model.Data{data=dto.TokenResponse}	"Token response consisting of access and refresh token"
//	@Failure		409				{object}	model.Error{error=usecase.AppError}
//	@Failure		422				{object}	model.Error{error=usecase.AppError}
//	@Failure		500				{object}	model.Error{error=usecase.AppError}
//	@Router			/credentials/signup [post]
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

// LoginHandler godoc
//
//	@Summary		Log in handler
//	@Description	Handles log in for new users
//	@Tags			Credentials
//	@Accept			json
//	@Produce		json
//	@Param			logInRequest	body		dto.LoginRequest					true	"Login request body"
//	@Success		200				{object}	model.Data{data=dto.TokenResponse}	"Token response consisting of access and refresh token"
//	@Failure		400				{object}	model.Error{error=usecase.AppError}
//	@Failure		401				{object}	model.Error{error=usecase.AppError}
//	@Failure		404				{object}	model.Error{error=usecase.AppError}
//	@Failure		422				{object}	model.Error{error=usecase.AppError}
//	@Failure		500				{object}	model.Error{error=usecase.AppError}
//	@Router			/credentials/login [post]
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

// RefreshHandler godoc
//
//	@Summary		Refresh handler
//	@Description	Handles requesting a new set of access and refresh token from the given previous refresh token
//	@Tags			Credentials
//	@Accept			json
//	@Produce		json
//	@Param			refreshRequest	body		dto.RefreshRequest					true	"Refresh token request body"
//	@Success		200				{object}	model.Data{data=dto.TokenResponse}	"Token response consisting of access and refresh token"
//	@Failure		400				{object}	model.Error{error=usecase.AppError}
//	@Failure		401				{object}	model.Error{error=usecase.AppError}
//	@Failure		404				{object}	model.Error{error=usecase.AppError}
//	@Failure		422				{object}	model.Error{error=usecase.AppError}
//	@Failure		500				{object}	model.Error{error=usecase.AppError}
//	@Router			/credentials/refresh [post]
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
