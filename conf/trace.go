package conf

type TraceStruct struct {
	Express string `json:"waybillCode"`
}

type TraceRoute struct {
	OpeTitle    string `json:"ope_title"`
	OpeTime     string `json:"ope_time"`
	OpeRemark   string `json:"ope_remark"`
	OpeName     string `json:"ope_name"`
	WaybillCode string `json:"waybill_code"`
}

type TraceApiDtosStruct struct {
	TraceApiDtos []TraceRoute `json:"trace_api_dtos"`
	Code         string  `json:"code"`
}

type TraceResStruct struct {
	Response TraceApiDtosStruct `json:"jingdong_etms_trace_get_responce"`
}
