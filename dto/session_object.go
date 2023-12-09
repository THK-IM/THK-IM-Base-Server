package dto

type (
	GetUploadParamsReq struct {
		SId      int64  `json:"s_id" form:"s_id"`
		UId      int64  `json:"u_id" form:"u_id"`
		ClientId int64  `json:"client_id" form:"client_id"`
		FName    string `json:"f_name" form:"f_name"`
	}

	GetUploadParamsRes struct {
		Id     int64             `json:"id" form:"id"`
		Engine string            `json:"engine" form:"engine"`
		Url    string            `json:"url" form:"url"`
		Method string            `json:"method" form:"method"`
		Params map[string]string `json:"params" form:"params"`
	}

	GetDownloadUrlReq struct {
		UId int64 `json:"u_id"`
		Id  int64 `json:"id" form:"id"`
	}
)
