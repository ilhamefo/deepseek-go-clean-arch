package constant

const (
	INVALID_REQUEST_BODY = "invalid_request_body"
	SUCCESS_EXPORT       = "export_rekap_success"
	VALIDATION_ERROR     = "validation_failed"
	ACCESS_TOKEN         = "access_token"
	REFRESH_TOKEN        = "refresh_token"
	SQL_ERROR            = "sql_error"

	ErrDecrypt              string = "Error while decrypting payload:"
	ErrGenerateFromPassword string = "Error generate from password:"
	ErrConvertToInt         string = "Error convert to int:"
	ErrGetInduk             string = "Error get induk :"
	ErrFormatDate           string = "Error format date: "
	ErrFormatTime           string = "Error convert time : "
	ErrForbidden            string = "forbidden"
	ErrInvalidFileFormat    string = "Invalid file format"
	ErrLevelFormat          string = "invalid level, must be 'ulp', 'up3', 'uid', 'reg', or 'nasional'"
	ErrUnitFilter           string = "Error get unit filter: "
	ErrGetLastUpdated       string = "Error get last updated: "
	ErrParallelTask         string = "Error run parallel task: "
	ErrGetData              string = "Error get data: "
)
