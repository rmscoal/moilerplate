package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto/mapper"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/model"
)

type UserProfileController struct {
	model.BaseControllerV1
	uc usecase.IUserProfileUseCase
}

func NewUserProfileController(rg *gin.RouterGroup, uc usecase.IUserProfileUseCase) {
	controller := new(UserProfileController)
	controller.uc = uc

	r := rg.Group("/profiles")
	{
		r.GET("/me", controller.getProfile)
		r.PUT("/email", controller.editEmailsHandler)
	}
}

/*
*************************************************
Controllers
*************************************************
*/
// ModifyUserEmails godoc
//
//	@Summary		Replaces user's emails
//	@Description	Handles for user to modify its emails
//	@Tags			protected,profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization		header		string										true	"Bearer + your access token"
//	@Param			modifyEmailRequest	body		dto.ModifyEmailRequest						true	"email changes request body"
//	@Success		200					{object}	model.Data{data=dto.ModifyEmailResponse}	"List of user's emails with its changes"
//	@Failure		400					{object}	model.Error{error=usecase.AppError}
//	@Failure		401					{object}	model.Error{error=usecase.AppError}
//	@Failure		404					{object}	model.Error{error=usecase.AppError}
//	@Failure		422					{object}	model.Error{error=usecase.AppError}
//	@Failure		500					{object}	model.Error{error=usecase.AppError}
//	@Router			/ptd/profiles/email [put]
func (controller *UserProfileController) editEmailsHandler(c *gin.Context) {
	var raw dto.ModifyEmailRequest
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		controller.ClientError(c, err)
		return
	}

	req := mapper.Profile.MapModifyEmailRequestToUserDomain(c.Keys["userId"].(string), raw)
	data, err := controller.uc.ModifyEmailAddress(c.Request.Context(), req)
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	res := mapper.Profile.MapUserDomainToModifyEmailResponse(data)
	controller.Ok(c, res)
}

// GetProfile godoc
//
//	@Summary		Get's the user's profile
//	@Description	Handles the retrieval of user's full profile data.
//	@Tags			protected,profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string										true	"Bearer + your access token"
//	@Success		200				{object}	model.Data{data=dto.FullProfileResponse}	"A full profile response of the user"
//	@Failure		400				{object}	model.Error{error=usecase.AppError}
//	@Failure		401				{object}	model.Error{error=usecase.AppError}
//	@Failure		404				{object}	model.Error{error=usecase.AppError}
//	@Failure		500				{object}	model.Error{error=usecase.AppError}
//	@Router			/ptd/profiles/me [get]
func (controller *UserProfileController) getProfile(c *gin.Context) {
	data, err := controller.uc.RetrieveProfile(c.Request.Context(), c.Keys["userId"].(string))
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Ok(c, mapper.Profile.MapUserDomainToFullProfileResponse(data))
}
