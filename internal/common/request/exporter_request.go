package request

type RekapRequest struct {
	UnitCode      string `json:"unit_code" form:"unit_code" validate:"" example:""`
	Area          string `json:"id_area" form:"id_area" validate:"" example:"52000"`
	Induk         string `json:"id_induk" form:"id_induk" validate:"" example:""`
	Pusat         string `json:"id_pusat" form:"id_pusat" validate:"" example:""`
	IsDBPlnMobile bool   `json:"is_db_plnmobile" form:"is_db_plnmobile" validate:"boolean" example:"false"`
	DateStart     string `json:"date_start" form:"date_start" validate:"required,datetime=2006/01/02,max=100" example:"2026/02/01"`
	DateEnd       string `json:"date_end" form:"date_end" validate:"required,datetime=2006/01/02,max=100" example:"2026/12/31"`
	Limit         int    `json:"limit" form:"limit" validate:"" example:"1000"`
	Offset        int    `json:"offset" form:"offset" validate:"" example:"0"`
}
