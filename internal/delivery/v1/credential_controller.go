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
//	@Param			signUpRequest	body		dto.SignUpRequest	true	"Signup request body"
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

	controller.Created(c, mapper.Credential.UserDomainToSignUpResponse(user))
}

// LoginHandler godoc
//
//	@Summary		Log in handler
//	@Description	Handles log in for signed up users
//	@Tags			Credentials
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		dto.LoginRequest	true	"Login request body"
//	@Success		200				{object}	model.Data{data=vo.Token}
//	@Failure		409				{object}	model.Error{error=usecase.AppError}
//	@Failure		422				{object}	model.Error{error=usecase.AppError}
//	@Failure		500				{object}	model.Error{error=usecase.AppError}
//	@Router			/credentials/login [post]
func (controller *CredentialController) loginHandler(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		controller.ClientError(c, usecase.NewClientError("Body", err))
		return
	}

	token, err := controller.uc.Login(c.Request.Context(), mapper.Credential.LoginRequestToUserDomain(req))
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Ok(c, token)
}

// RefreshAccessHandler godoc
//
//	@Summary		Refresh access handler
//	@Description	Handles log in for refresh users
//	@Tags			Credentials
//	@Accept			json
//	@Produce		json
//	@Param			refreshRequest	body		dto.RefreshRequest	true	"refresh request body"
//	@Success		200				{object}	model.Data{data=vo.Token}
//	@Failure		409				{object}	model.Error{error=usecase.AppError}
//	@Failure		422				{object}	model.Error{error=usecase.AppError}
//	@Failure		500				{object}	model.Error{error=usecase.AppError}
//	@Router			/credentials/refresh [post]
func (controller *CredentialController) refreshHandler(c *gin.Context) {
	var req dto.RefreshRequest

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		controller.ClientError(c, usecase.NewClientError("Body", err))
		return
	}

	token, err := controller.uc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Ok(c, token)
}
