package main

type RequestHead struct {
	RequestId string `json:"requestId"`
	UserType  int32  `json:"userType"`
	Cmd       int32  `json:"cmd"`
	WsId      string `json:"wsId"`
}

type ResponseHead struct {
	RequestId string `json:"requestId"`
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Cmd       int32  `json:"cmd"`
}

type wholeSalerRegisterReqData struct {
	WsName    string `json:"wsName"`
	WsCompany string `json:"wsCompany"`
	WsMobile  string `json:"wsMobile"`
}

type wholeSalerRegisterReq struct {
	RequestHead
	OpenId string                    `json:"openId"`
	UserId string                    `json:"userId"`
	Data   wholeSalerRegisterReqData `json:"data"`
}

type wholeSalerRegisterRespData struct {
	WsId           string `json:"wsId"`
	WsName         string `json:"wsName"`
	WsCompany      string `json:"wsCompany"`
	WsMobile       string `json:"wsMobile"`
	WsIdentityCode string `json:"wsIdentityCode"`
}

type wholeSalerRegisterResp struct {
	ResponseHead
	Data wholeSalerRegisterRespData `json:"data"`
}

type userLoginReq struct {
	RequestId string `json:"requestId"`
	Passwd    string `json:"passwd"`
	UserType  int32  `json:"userType"`
	UserName  string `json:"userName"`
}

type userLoginResp struct {
	RequestId string `json:"requestId"`
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	OpenId    string `json:"openId"`
	UserId    string `json:"userId"`
	UserType  int32  `json:"userType"`
	UserName  string `json:"userName"`
	HeadIco   string `json:"headIco"`
}

type QueryUserReqData struct {
	SpId   string `json:"spId"`
	WxCode string `json:"wxCode"`
}

type QueryUserReq struct {
	RequestHead
	Data QueryUserReqData `json:"data"`
}

type QueryUserRespData struct {
	OpenId   string `json:"openId"`
	UserId   string `json:"userId"`
	UserType int32  `json:"userType"`
	UserName string `json:"userName"`
	HeadIco  string `json:"headIco"`
}

type QueryUserResp struct {
	ResponseHead ResponseHead
	Data         QueryUserRespData `json:"data"`
}
