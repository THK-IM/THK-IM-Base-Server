package errorx

var (
	ErrParamsError         = NewErrorX(4000000, "Params Error")
	ErrPermission          = NewErrorX(4000001, "Permission denied")
	ErrNotSupportReCommit  = NewErrorX(4000002, "Not support Recommit")
	ErrNotFound            = NewErrorX(4000003, "Not Found")
	ErrInternalServerError = NewErrorX(5000000, "Internal Server err")
	ErrServerBusy          = NewErrorX(5000001, "Server busy")
)
