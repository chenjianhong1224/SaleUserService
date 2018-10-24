package main

import (
	"fmt"

	"go.uber.org/zap"
)

type wholesaler_service struct {
	d *dbOperator
}

func (m *wholesaler_service) queryWholesaler(mobile string, company string) (*TWholeSaler, error) {
	args := []interface{}{}
	args = append(args, company)
	args = append(args, mobile)
	tmp := TWholeSaler{}
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
	return queryRep.Rows[0].(*TWholeSaler), nil
}
