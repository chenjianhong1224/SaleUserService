package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type wholesaler_service struct {
	d *dbOperator
}

func (m *wholesaler_service) makepw(pwl int) string {
	chars := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
	clen := float64(len(chars))
	res := ""
	rand.Seed(time.Now().Unix())
	for i := 0; i < pwl; i++ {
		rfi := int(clen * rand.Float64())
		res += fmt.Sprintf("%c", chars[rfi])
	}
	return res
}

func (m *wholesaler_service) addWholesaler(req wholeSalerRegisterReq) (int64, string, error) {
	passwd := m.makepw(8)
	args2 := []interface{}{}
	uid, _ := uuid.NewV4()
	args2 = append(args2, uid.String())
	args2 = append(args2, req.WsName)
	data := []byte(passwd)
	has := md5.Sum(data)
	args2 = append(args2, fmt.Sprintf("%x", has))
	args2 = append(args2, uid.String())
	args2 = append(args2, req.UserId)
	args2 = append(args2, req.UserId)
	execReq2 := SqlExecRequest{
		SQL:  "insert into t_user(user_uuid, user_name, passwd, agent_uuid, user_type, user_status, create_time, create_user, update_time, update_user) values(?,?,?,?,1,1,now(),?,now(),?)",
		Args: args2,
	}
	args1 := []interface{}{}
	args1 = append(args1, uid.String())
	args1 = append(args1, req.WsMobile)
	args1 = append(args1, req.WsCompany)
	args1 = append(args1, req.WsMobile)
	args1 = append(args1, req.UserId)
	args1 = append(args1, req.UserId)
	execReq1 := SqlExecRequest{
		SQL:  "insert into t_wholesaler(saler_uuid, saler_name, company, mobile, saler_status, create_time, create_user, update_time, update_user) values(?,?,?,?,1,now(),?,now(),?)",
		Args: args1,
	}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		tWholeSaler, err := m.queryWholesaler(req.WsMobile, req.WsCompany)
		if err == nil {
			return tWholeSaler.Saler_id, passwd, nil
		}
	}
	zap.L().Error(fmt.Sprintf("add wholesaler[%s,%s] error:%s", req.WsCompany, req.WsMobile, err.Error()))
	return -1, "", err
}

func (m *wholesaler_service) queryWholesaler(mobile string, company string) (*tWholeSaler, error) {
	args := []interface{}{}
	args = append(args, company)
	args = append(args, mobile)
	tmp := tWholeSaler{}
	queryReq := &SqlQueryRequest{
		SQL:         "select saler_id, saler_uuid, saler_name, company, mobile , saler_status , create_time, create_user, update_time, update_user, remark from t_wholesaler where company = ? and mobile = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query wholesaler[%s,%s] error:%s", company, mobile, queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	return queryRep.Rows[0].(*tWholeSaler), nil
}
