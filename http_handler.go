package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type clientInfo struct {
	ipStr string
	ipNum int32
	port  int32
}

type httpHandler struct {
	cfg               *Config
	wholesalersv      *wholesaler_service
	usersv            *user_service
	wxpluginprogramsv *wxplugin_program_service
}

func (ci *clientInfo) inetAton() {
	ip := net.ParseIP(ci.ipStr)
	ci.ipNum = int32(binary.BigEndian.Uint32(ip.To4()))
}

func (m *httpHandler) start() error {
	//start http server
	s := &http.Server{
		Addr:           m.cfg.Server.Endpoint,
		Handler:        nil,
		ReadTimeout:    m.cfg.Server.HttpReadTimeout,
		WriteTimeout:   m.cfg.Server.HttpWriteTimeout,
		MaxHeaderBytes: int(m.cfg.Server.MaxHeadSize),
	}
	http.HandleFunc("/api", m.process)
	go s.ListenAndServe()

	return nil
}

func (m *httpHandler) ivalidResp(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (m *httpHandler) getClientInfo(r *http.Request) *clientInfo {
	cliIp, cliPort, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("userip: %q is not IP:port", r.RemoteAddr))
		return &clientInfo{ipNum: 0, port: 0}
	} else {
		zap.L().Debug(fmt.Sprintf("package from %s:%s", cliIp, cliPort))
		p, e := strconv.Atoi(cliPort)
		if e != nil {
			zap.L().Error(fmt.Sprintf("strconv Atoi port fail"))
			p = 0
		}

		ci := &clientInfo{
			ipStr: cliIp,
			port:  int32(p),
			ipNum: 0,
		}

		ci.inetAton()
		return ci
	}
}

func (m *httpHandler) process(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		zap.L().Info(fmt.Sprintf("get method not support, method:%s", r.Method))
		statObj.statHandler.StatCount(StatInvalidMethodReq)
		m.ivalidResp(w)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		statObj.statHandler.StatCount(StatReadBody)
		m.ivalidResp(w)
		return
	} else {
		zap.L().Debug(fmt.Sprintf("recv body len:%d content:%s", len(body), body))
		var req RequestHead
		err := json.Unmarshal(body, &req)
		if err != nil {
			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
			m.ivalidResp(w)
			return
		}
		if req.Cmd == 130 { //用户信息查询接口
			m.queryUser(body, w)
		} else if req.Cmd == 132 { //用户登录接口
			m.userLogin(body, w)
		} else if req.Cmd == 134 { //用户退出接口
			m.userLogout(body, w)
		} else if req.Cmd == 136 { //用户修改密码接口
			m.changePasswd(body, w)
		} else {
			var respHead ResponseHead
			respHead = ResponseHead{RequestId: req.RequestId, ErrorCode: 9999, Cmd: req.Cmd, ErrorMsg: "cmd不合法"}
			jsonData, err := json.Marshal(respHead)
			if err != nil {
				zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
				m.ivalidResp(w)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(jsonData))
			return
		}
	}
}

func (m *httpHandler) userLogin(body []byte, w http.ResponseWriter) {
	var req UserLoginReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	errMsg, openId, _, err := m.getWxUserInfo(req.Data.SpId, req.Data.WxCode)
	if err != nil {
		zap.L().Error(fmt.Sprintf("userLogin error %s", err.Error()))
		m.ivalidResp(w)
		return
	} else if errMsg != "" {
		var respHead ResponseHead
		respHead = ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 9999, ErrorMsg: errMsg, Cmd: 133}
		jsonData, err := json.Marshal(respHead)
		if err != nil {
			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
			m.ivalidResp(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
		return
	}
	var resp UserLoginResp
	if req.UserType != 1 { //零售商登录不需要密码，其他都需要
		tUser, err := m.usersv.queryUserByPasswd(req.Data.Passwd, req.Data.LoginName, req.UserType)
		if err == nil {
			if tUser == nil {
				resp = UserLoginResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 1, ErrorMsg: "用户" + req.Data.LoginName + "登录失败: 用户名或密码不正确", Cmd: 133}, UserLoginRespData{}}
			} else {
				err = m.usersv.bindUser(openId, tUser.User_uuid)
				if err == nil {
					err = m.usersv.login(tUser.User_uuid)
				}
				resp = UserLoginResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 0, Cmd: 133}, UserLoginRespData{OpenId: openId,
					UserId:   tUser.User_uuid,
					UserType: req.UserType,
					UserName: tUser.User_name.String,
					HeadIco:  ""}}

			}
		}
	} else {
		tUser, err := m.usersv.queryUserByOpenIdType(openId, req.UserType)
		if err == nil {
			var usrUUid string
			if tUser == nil {
				usrUUid, err = m.usersv.addRetailer(openId, req.WsId, "")
				if err != nil {
					zap.L().Error(fmt.Sprintf("login add addRetailer error %s", err.Error()))
					m.ivalidResp(w)
					return
				}
			} else {
				usrUUid = tUser.User_uuid
			}
			err = m.usersv.login(usrUUid)
			resp = UserLoginResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 0, Cmd: 133}, UserLoginRespData{OpenId: openId,
				UserId:   usrUUid,
				UserType: req.UserType,
				UserName: "",
				HeadIco:  ""}}
		}
	}
	if err != nil {
		zap.L().Error(fmt.Sprintf("doUserLogin error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) userLogout(body []byte, w http.ResponseWriter) {
	var req UserLogoutReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	zap.L().Info(fmt.Sprintf("request %+v\n", req))
	err = m.usersv.logout(req.UserId)
	var resp UserLogoutResp
	if err != nil {
		zap.L().Error(fmt.Sprintf("logout error %s", err.Error()))
		resp = UserLogoutResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 9999, Cmd: 135, ErrorMsg: "登出失败:" + err.Error()}}
	} else {
		resp = UserLogoutResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 0, Cmd: 135}}
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) changePasswd(body []byte, w http.ResponseWriter) {
	var req ChangePasswdReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp ChangePasswdResp
	if req.UserType == 0 || req.UserType == 2 {
		tUser, err := m.usersv.queryUserByUUidPasswd(req.data.OldPasswd, req.UserId, req.UserType)
		if err == nil && tUser != nil {
			err = m.usersv.changePasswd(req.data.Passwd, req.UserId)
			if err == nil {
				resp = ChangePasswdResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 0, Cmd: 137}}
			} else {
				resp = ChangePasswdResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 9999, Cmd: 137, ErrorMsg: "修改密码失败"}}
			}
		} else {
			zap.L().Info(fmt.Sprintf("修改密码时找不到对应用户 %s, %s, %d", req.UserId, req.data.OldPasswd, req.UserType))
			resp = ChangePasswdResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 9999, Cmd: 137, ErrorMsg: "修改密码失败"}}
		}
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) getWxUserInfo(spId string, wxCode string) (errMsg string, openId string, sessionKey string, err error) {
	var appid string
	var secret string
	tWxPluginProgram, err := m.wxpluginprogramsv.queryWxpluginProgram(spId)
	if err != nil {
		zap.L().Error(fmt.Sprintf("getWxUserInfo error %s", err.Error()))
		return "", "", "", err
	} else if tWxPluginProgram == nil {
		return "系统未配置对应的小程序", "", "", nil
	}
	appid = tWxPluginProgram.Appid
	secret = tWxPluginProgram.Appsecrete
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + wxCode + "&grant_type=authorization_code"
	zap.L().Debug("tencent authorization url = " + url)
	resp, err := http.Get(url)
	if err != nil {
		zap.L().Error(fmt.Sprintf("get wx session_key error %s", err.Error()))
		return "", "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error(fmt.Sprintf("get wx session_key error %s", err.Error()))
		return "", "", "", err
	}
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dat); err == nil {
		zap.L().Debug("tencent return = " + string(body))
		openid := dat["openid"]
		session_key := dat["session_key"]
		zap.L().Debug(fmt.Sprintf("openid:%s", openid.(string)))
		zap.L().Debug(fmt.Sprintf("session_key:%s", session_key.(string)))
		return "", openid.(string), session_key.(string), err
	}
	return "", "", "", err
}

func (m *httpHandler) queryUser(body []byte, w http.ResponseWriter) {
	var req QueryUserReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	zap.L().Info(fmt.Sprintf("request %+v\n", req))
	errMsg, openid, _, err := m.getWxUserInfo(req.Data.SpId, req.Data.WxCode)
	if err != nil {
		zap.L().Error(fmt.Sprintf("queryUser error %s", err.Error()))
		m.ivalidResp(w)
		return
	} else if errMsg != "" {
		var respHead ResponseHead
		respHead = ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 9999, ErrorMsg: errMsg, Cmd: 131}
		jsonData, err := json.Marshal(respHead)
		if err != nil {
			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
			m.ivalidResp(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
		return
	}
	tUsers, err := m.usersv.queryUserByOpenId(openid)
	var resp QueryUserResp
	if err == nil {
		var data QueryUserRespData
		if len(tUsers) != 0 {
			var findUser bool
			findUser = false
			for i := 0; i < len(tUsers); i++ {
				findUser = true
				var respHead ResponseHead
				respHead = ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 0, Cmd: 131}
				data = QueryUserRespData{OpenId: tUsers[i].Open_id.String, UserId: tUsers[i].User_uuid, UserType: tUsers[i].User_type, UserName: tUsers[i].User_name.String, HeadIco: tUsers[i].Head_portrait.String}
				resp = QueryUserResp{ResponseHead: respHead, Data: data}
				break
			}
			if findUser == false {
				resp = QueryUserResp{ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 1, Cmd: 131, ErrorMsg: "未查到对应的用户"}, QueryUserRespData{}}
			}
		} else {
			resp = QueryUserResp{ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 1, Cmd: 131, ErrorMsg: "未查到对应的用户"}, QueryUserRespData{}}
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
			m.ivalidResp(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
		return
	}
	zap.L().Error(fmt.Sprintf("queryUser error %s", err.Error()))
	m.ivalidResp(w)
}
