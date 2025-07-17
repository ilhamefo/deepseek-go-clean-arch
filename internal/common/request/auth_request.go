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
type SearchRequest struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"required" example:"induk@gmail.com"`
}

type GetUnitRequest struct {
	Level string `json:"level" query:"level" form:"level" validate:"required,oneof=0 1 2 3" example:"induk@gmail.com"`
}

type UpdateUserRequest struct {
	ID                   string   `json:"id" params:"id" validate:"required,max=100" example:"2"`
	Email                string   `json:"email" validate:"required,max=100" example:"staff@gmail.com"`
	Username             string   `json:"username" validate:"required,max=100" example:"Staff"`
	FullName             string   `json:"full_name" validate:"required,max=100" example:"Staff Name"`
	Level                string   `json:"level" validate:"numeric,max=100" example:"3"`
	Jabatan              string   `json:"jabatan" validate:"max=100" example:"Staff Jabatan"`
	NIP                  string   `json:"nip" validate:"max=100" example:"1234567890"`
	UnitCode             string   `json:"unit_code" validate:"max=100" example:"52001"`
	UnitName             string   `json:"unit_name" validate:"max=100" example:"Unit Name"`
	Status               string   `json:"status" validate:"numeric,max=100" example:"1"`
	Roles                []string `json:"roles" validate:"required,dive,max=100" example:"3"`
	Password             string   `json:"password" validate:""`
	PasswordConfirmation string   `json:"password_confirmation" binding:"required_with=Password,eqfield=Password"`
}
