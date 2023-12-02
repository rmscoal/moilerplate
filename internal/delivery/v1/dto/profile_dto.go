package dto

type FullProfileResponse struct {
	FirstName   string        `json:"firstName" example:"firstName"`
	LastName    string        `json:"lastName" example:"lastName"`
	Username    string        `json:"username" example:"username"`
	PhoneNumber string        `json:"phoneNumber" example:"0823456786543"`
	Emails      []EmailDetail `json:"emails"`
}

type ModifyEmailRequest struct {
	Emails []EmailDetail `json:"emails"`
}

type EmailDetail struct {
	Email     string `json:"email,omitempty" example:"email@email.com"`
	IsPrimary bool   `json:"isPrimary"`
}

type ModifyEmailResponse struct {
	UserId string        `json:"userId,omitempty" example:"userID"`
	Emails []EmailDetail `json:"emails,omitempty"`
}
