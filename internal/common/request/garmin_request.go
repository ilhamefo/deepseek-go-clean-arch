package request

type RefreshActivitiesRequest struct {
	Token   string `json:"token" query:"token" form:"token" validate:"required"`
	Cookies string `json:"cookies" query:"cookies" form:"cookies" validate:"required"`
}

type ActivityRequest struct {
	ActivityID string `json:"activity_id" query:"activity_id" form:"activity_id" validate:"required"`
}
