package request

type GoogleCallbackRequest struct {
	Code        string `json:"code" query:"code" validate:"required"`
	State       string `json:"state" query:"state" validate:"required"`
	StateCookie string
}
