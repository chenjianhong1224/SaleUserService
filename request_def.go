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

type wholeSalerRegisterReq struct {
	RequestId string `json:"requestId"`
	OpenId    string `json:"openId"`
	UserId    string `json:"userId"`
	UserType  int32  `json:"userType"`
	WsName    string `json:"wsName"`
	WsCompany string `json:"wsCompany"`
	WsMobile  string `json:"wsMobile"`
}

type wholeSalerRegisterResp struct {
	RequestId      string `json:"requestId"`
	ErrorCode      int32  `json:"errorCode"`
	ErrorMsg       string `json:"errorMsg"`
	WsId           string `json:"wsId"`
	UserType       int32  `json:"userType"`
	WsName         string `json:"wsName"`
	WsCompany      string `json:"wsCompany"`
	WsMobile       string `json:"wsMobile"`
	WsIdentityCode string `json:"wsIdentityCode"`
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
	RequestHead RequestHead
	Data        QueryUserReqData `json:"data"`
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
