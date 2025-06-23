package request

type RekapRequest struct {
	UnitCode  string `json:"unit_code" form:"unit_code" validate:"" example:""`
	Area      string `json:"id_area" form:"id_area" validate:"" example:"52000"`
	Induk     string `json:"id_induk" form:"id_induk" validate:"" example:""`
	Pusat     string `json:"id_pusat" form:"id_pusat" validate:"" example:""`
	DateStart string `json:"date_start" form:"date_start" validate:"required,datetime=2006/01/02,max=100" example:"2025/05/01"`
	DateEnd   string `json:"date_end" form:"date_end" validate:"required,datetime=2006/01/02,max=100" example:"2025/05/27"`
}
