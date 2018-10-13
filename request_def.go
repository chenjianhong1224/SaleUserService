package main

type RequestHead struct {
	RequestId string `json:"requestId"`
	UserType  int32  `json:"userType"`
	Cmd       int32  `json:"cmd"`
	WsId      string `json:"wsId"` //uuid
}

type ResponseHead struct {
	RequestId string `json:"requestId"`
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Cmd       int32  `json:"cmd"`
}

type WholeSalerRegisterReqData struct {
	WsName    string `json:"wsName"`
	WsCompany string `json:"wsCompany"`
	WsMobile  string `json:"wsMobile"`
}

type WholeSalerRegisterReq struct {
	RequestHead
	OpenId string                    `json:"openId"`
	UserId string                    `json:"userId"`
	Data   WholeSalerRegisterReqData `json:"data"`
}

type WholeSalerRegisterRespData struct {
	WsId           string `json:"wsId"`
	WsName         string `json:"wsName"`
	WsCompany      string `json:"wsCompany"`
	WsMobile       string `json:"wsMobile"`
	WsIdentityCode string `json:"wsIdentityCode"`
}

type WholeSalerRegisterResp struct {
	ResponseHead
	Data WholeSalerRegisterRespData `json:"data"`
}

type UserLoginReq struct {
	RequestHead
	Data UserLoginReqData `json:"data"`
}

type UserLoginReqData struct {
	Passwd    string `json:"passwd"`
	SpId      string `json:"spId"`
	WxCode    string `json:"wxCode"`
	LoginName string `json:"loginName"`
}

type UserLoginResp struct {
	ResponseHead
	Data UserLoginRespData `json:"data"`
}

type UserLoginRespData struct {
	OpenId   string `json:"openId"`
	UserId   string `json:"userId"`
	UserType int32  `json:"userType"`
	UserName string `json:"userName"`
	HeadIco  string `json:"headIco"`
}

type UserLogoutReq struct {
	RequestId string `json:"requestId"`
	UserType  int32  `json:"userType"`
	Cmd       int32  `json:"cmd"`
	UserId    string `json:"userId"`
	wsId      string `json:"wsId"`
}

type UserLogoutResp struct {
	ResponseHead
}

type ChangePasswdReq struct {
	RequestId string              `json:"requestId"`
	UserType  int32               `json:"userType"`
	Cmd       int32               `json:"cmd"`
	UserId    string              `json:"userId"`
	wsId      string              `json:"wsId"`
	data      ChangePasswdReqData `json:"data"`
}

type ChangePasswdReqData struct {
	OldPasswd string `json:"oldPasswd"`
	Passwd    string `json:"passwd"`
}

type ChangePasswdResp struct {
	ResponseHead ResponseHead
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
