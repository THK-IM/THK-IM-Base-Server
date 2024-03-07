package errorx

var (
	ErrParamsError         = NewErrorX(4000000, "Params Error")
	ErrNotFound            = NewErrorX(4000001, "Not Found")
	ErrNotSupportReCommit  = NewErrorX(4000002, "Not support Recommit")
	ErrPermission          = NewErrorX(4000003, "Permission denied")
	ErrInternalServerError = NewErrorX(5000000, "Internal Server err")
	ErrServerBusy          = NewErrorX(5000001, "Server busy")
)
