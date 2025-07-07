package request

type GoogleCallbackRequest struct {
	Code        string `json:"code" query:"code" form:"code" validate:"required"`
	State       string `json:"state" query:"state" form:"state" validate:"required"`
	StateCookie string
}

type LoginRequest struct {
	Email    string `json:"email" query:"email" form:"email" validate:"required,email" example:"ilham@oninyon.com"`
	Password string `json:"password" query:"password" form:"password" validate:"required" example:"password"`
}
