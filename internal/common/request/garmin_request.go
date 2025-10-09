package request

type GarminBasicRequest struct {
	Cookies         string `json:"cookies" query:"cookies" form:"cookies" validate:"required"`
	GarminCsrfToken string `json:"garmin_csrf_token" query:"garmin_csrf_token" form:"garmin_csrf_token" validate:"required"`
}

type ActivityRequest struct {
	ActivityID string `json:"activity_id" query:"activity_id" form:"activity_id" validate:"required"`
}

type GarminByDateRequest struct {
	GarminBasicRequest
	Date string `json:"date" form:"date" validate:"required,datetime=2006-01-02,max=100" example:"2025-10-05"`
}
