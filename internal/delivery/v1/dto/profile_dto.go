package dto

type FullProfileResponse struct {
	FirstName   string        `json:"firstName"`
	LastName    string        `json:"lastName"`
	Username    string        `json:"username"`
	PhoneNumber string        `json:"phoneNumber"`
	Emails      []EmailDetail `json:"emails"`
}

type ModifyEmailRequest struct {
	Emails []EmailDetail `json:"emails"`
}

type EmailDetail struct {
	Email     string `json:"email,omitempty"`
	IsPrimary bool   `json:"isPrimary"`
}

type ModifyEmailResponse struct {
	UserId string        `json:"userId,omitempty"`
	Emails []EmailDetail `json:"emails,omitempty"`
}
