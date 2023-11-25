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

func (controller *UserProfileController) getProfile(c *gin.Context) {
	data, err := controller.uc.RetrieveProfile(c.Request.Context(), c.Keys["userId"].(string))
	if err != nil {
		controller.SummariesUseCaseError(c, err)
		return
	}

	controller.Ok(c, mapper.Profile.MapUserDomainToFullProfileResponse(data))
}
