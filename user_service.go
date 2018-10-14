package main

import (
	"crypto/md5"
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type user_service struct {
	d *dbOperator
}

func (m *user_service) queryUser(openId string, user_uuid string) (*tUser, error) {
	args := []interface{}{}
	args = append(args, openId)
	args = append(args, user_uuid)
	tmp := tUser{}
	queryReq := &SqlQueryRequest{
		SQL:         "select User_id, User_uuid, User_name, Passwd, Open_id, Other_from, Nickname, Head_portrait, Agent_uuid, User_type, User_status, User_token, Expiry_time, Create_time, Create_user, Update_time, Update_user, Remark from t_user where open_id = ? and user_uuid = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query user[%s,%s] error:%s", openId, user_uuid, queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	return queryRep.Rows[0].(*tUser), nil
}

func (m *user_service) queryUserByUUidPasswd(passwd string, user_uuid string, user_type int32) (*tUser, error) {
	args := []interface{}{}
	data := []byte(passwd)
	has := md5.Sum(data)
	args = append(args, fmt.Sprintf("%x", has))
	args = append(args, user_uuid)
	args = append(args, user_type)
	tmp := tUser{}
	queryReq := &SqlQueryRequest{
		SQL:         "select User_id, User_uuid, User_name, Passwd, Open_id, Other_from, Nickname, Head_portrait, Agent_uuid, User_type, User_status, User_token, Expiry_time, Create_time, Create_user, Update_time, Update_user, Remark from t_user where Passwd = ? and User_uuid = ? and User_type = ? and User_status = 1",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query user[%s,%d] error:%s", user_uuid, user_type, queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	return queryRep.Rows[0].(*tUser), nil
}

func (m *user_service) queryUserByPasswd(passwd string, login_name string, user_type int32) (*tUser, error) {
	args := []interface{}{}
	data := []byte(passwd)
	has := md5.Sum(data)
	args = append(args, fmt.Sprintf("%x", has))
	args = append(args, login_name)
	args = append(args, user_type)
	tmp := tUser{}
	queryReq := &SqlQueryRequest{
		SQL:         "select User_id, User_uuid, User_name, Passwd, Open_id, Other_from, Nickname, Head_portrait, Agent_uuid, User_type, User_status, User_token, Expiry_time, Create_time, Create_user, Update_time, Update_user, Remark from t_user where Passwd = ? and User_name = ? and User_type = ? and User_status = 1",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query user[%s,%d] error:%s", login_name, user_type, queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	return queryRep.Rows[0].(*tUser), nil
}

func (m *user_service) queryUserByOpenId(openId string) ([]tUser, error) {
	args := []interface{}{}
	args = append(args, openId)
	tmp := tUser{}
	queryReq := &SqlQueryRequest{
		SQL:         "select User_id, User_uuid, User_name, Passwd, Open_id, Other_from, Nickname, Head_portrait, Agent_uuid, User_type, User_status, User_token, Expiry_time, Create_time, Create_user, Update_time, Update_user, Remark from t_user where open_id = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query user[%s] error:%s", openId, queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	tUsers := []tUser{}
	for i := 0; i < len(queryRep.Rows); i++ {
		tUsers = append(tUsers, queryRep.Rows[i].(tUser))
	}
	return tUsers, nil
}

func (m *user_service) queryUserByOpenIdType(openId string, userType int32) (*tUser, error) {
	args := []interface{}{}
	args = append(args, openId)
	args = append(args, userType)
	tmp := tUser{}
	queryReq := &SqlQueryRequest{
		SQL:         "select User_id, User_uuid, User_name, Passwd, Open_id, Other_from, Nickname, Head_portrait, Agent_uuid, User_type, User_status, User_token, Expiry_time, Create_time, Create_user, Update_time, Update_user, Remark from t_user where open_id = ? and User_type = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query user[%s] error:%s", openId, queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	return queryRep.Rows[0].(*tUser), nil
}

func (m *user_service) bindUser(openId string, user_uuid string) error {
	args1 := []interface{}{}
	args1 = append(args1, openId)
	args1 = append(args1, user_uuid)
	args1 = append(args1, user_uuid)
	execReq1 := &SqlExecRequest{
		SQL:  "update t_user set open_id = ?, update_user = ? , update_time = now()  where user_uuid = ?",
		Args: args1,
	}
	reply := m.d.dbCli.Query(execReq1)
	execRep, _ := reply.(*SqlExecReply)
	if execRep.Err != nil {
		zap.L().Error(fmt.Sprintf("bindUser user[%s,%s] error:%s", openId, user_uuid, execRep.Err.Error()))
		return execRep.Err
	}
	return nil
}

func (m *user_service) login(user_uuid string) error {
	args1 := []interface{}{}
	args1 = append(args1, 2)
	args1 = append(args1, user_uuid)
	args1 = append(args1, user_uuid)
	execReq1 := &SqlExecRequest{
		SQL:  "update t_user set user_status = ?, login_time = now()  where user_uuid = ?",
		Args: args1,
	}
	reply := m.d.dbCli.Query(execReq1)
	execRep, _ := reply.(*SqlExecReply)
	if execRep.Err != nil {
		zap.L().Error(fmt.Sprintf("login user[%s] error:%s", user_uuid, execRep.Err.Error()))
		return execRep.Err
	}
	return nil
}

func (m *user_service) logout(user_uuid string) error {
	args1 := []interface{}{}
	args1 = append(args1, 3)
	args1 = append(args1, user_uuid)
	execReq1 := &SqlExecRequest{
		SQL:  "update t_user set user_status = ?, logout_time = now()  where user_uuid = ?",
		Args: args1,
	}
	reply := m.d.dbCli.Query(execReq1)
	execRep, _ := reply.(*SqlExecReply)
	if execRep.Err != nil {
		zap.L().Error(fmt.Sprintf("logout user[%s] error:%s", user_uuid, execRep.Err.Error()))
		return execRep.Err
	}
	return nil
}

func (m *user_service) changePasswd(passwd string, user_uuid string) error {
	args := []interface{}{}
	data := []byte(passwd)
	has := md5.Sum(data)
	args = append(args, fmt.Sprintf("%x", has))
	args = append(args, user_uuid)
	execReq := &SqlExecRequest{
		SQL:  "update t_user set passwd = ?, update_time=now(), update_user=? where user_uuid = ?",
		Args: args,
	}
	reply := m.d.dbCli.Query(execReq)
	execRep, _ := reply.(*SqlExecReply)
	if execRep.Err != nil {
		zap.L().Error(fmt.Sprintf("changePasswd user[%s] error:%s", user_uuid, execRep.Err.Error()))
		return execRep.Err
	}
	return nil
}

func (m *user_service) addRetailer(openId string) (string, error) {
	uid, _ := uuid.NewV4()
	args2 := []interface{}{}
	args2 = append(args2, uid.String())
	args2 = append(args2, openId)
	args2 = append(args2, uid.String())
	args2 = append(args2, uid.String())
	execReq2 := &SqlExecRequest{
		SQL:  "insert into t_user(user_uuid, open_id, user_type, user_status, create_time, create_user, update_time, update_user) values(?,?,1,1,now(),?,now(),?)",
		Args: args2,
	}
	reply := m.d.dbCli.Query(execReq2)
	execRep, _ := reply.(*SqlExecReply)
	if execRep.Err != nil {
		zap.L().Error(fmt.Sprintf("addRetailer user[%s] error:%s", openId, execRep.Err.Error()))
		return "", execRep.Err
	}
	return uid.String(), nil
}
