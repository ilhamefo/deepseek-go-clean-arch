package request

import "encoding/json"

type ActivityDashboardRequest struct {
	Cursor int64       `json:"cursor" query:"cursor" form:"cursor" validate:""`
	Limit  int         `json:"limit" query:"limit" form:"limit" validate:"max=100"`
	Search string      `json:"search" query:"search" form:"search" validate:""`
	Type   json.Number `json:"type" query:"type" form:"type" validate:"numeric,max=100,min=0"`
}
