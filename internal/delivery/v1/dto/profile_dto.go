package dto

type ModifyEmailRequest struct {
	Emails []ModifyEmailDetailRequest `json:"emails"`
}

type ModifyEmailDetailRequest struct {
	Email     string `json:"email,omitempty"`
	IsPrimary bool   `json:"isPrimary"`
}

type ModifyEmailResponse struct {
	UserId string                     `json:"userId,omitempty"`
	Emails []ModifyEmailDetailRequest `json:"emails,omitempty"`
}
