package errorx

var (
	ErrParamsError        = NewErrorX(4000000, "Params Error")
	ErrPermission         = NewErrorX(4000001, "Permission denied")
	ErrNotSupportReCommit = NewErrorX(4000002, "Not support Recommit")
	ErrServerUnknown      = NewErrorX(5000000, "Server unknown err")
	ErrServerBusy         = NewErrorX(5000001, "Server busy")
)
