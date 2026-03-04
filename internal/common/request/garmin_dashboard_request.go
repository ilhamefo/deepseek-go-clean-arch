package request

type ActivityDashboardRequest struct {
	Cursor int64  `json:"cursor" query:"cursor" form:"cursor" validate:""`
	Limit  int    `json:"limit" query:"limit" form:"limit" validate:"max=100"`
	Search string `json:"search" query:"search" form:"search" validate:""`
	Type   *int   `json:"type" query:"type" form:"type" validate:"omitempty,max=100,min=0"`
}

type ActivityDetailsDashboardRequest struct {
	ActivityID int64 `json:"activity_id" query:"activity_id" form:"activity_id" validate:"required"`
}
