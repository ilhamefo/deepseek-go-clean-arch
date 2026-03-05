package request

type ActivityDashboardRequest struct {
	Cursor    int64  `json:"cursor" query:"cursor" form:"cursor" validate:""`
	Limit     int    `json:"limit" query:"limit" form:"limit" validate:"max=100"`
	Search    string `json:"search" query:"search" form:"search" validate:""`
	Type      *int   `json:"type" query:"type" form:"type" validate:"omitempty,max=100,min=0"`
	SortBy    string `json:"sort_by" query:"sort_by" form:"sort_by" validate:"omitempty,oneof=maxHr calories date name distance duration avgPace" example:"date"`
	SortOrder string `json:"sort_order" query:"sort_order" form:"sort_order" validate:"omitempty,oneof=asc desc" example:"desc"`
}

type ActivityDetailsDashboardRequest struct {
	ActivityID int64 `json:"activity_id" query:"activity_id" form:"activity_id" validate:"required"`
}
