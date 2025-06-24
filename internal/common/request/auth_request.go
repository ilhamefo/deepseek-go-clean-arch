package request

type GoogleCallbackRequest struct {
	Code        string `json:"code" query:"code" form:"code" validate:"required"`
	State       string `json:"state" query:"state" form:"state" validate:"required"`
	StateCookie string
}

type LoginRequest struct {
	Email    string `json:"email" query:"email" validate:"required"`
	Password string `json:"password" query:"password" validate:"required"`
}
