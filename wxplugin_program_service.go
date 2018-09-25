package main

import (
	"fmt"

	"go.uber.org/zap"
)

type wxplugin_program_service struct {
	d *dbOperator
}

func (m *wxplugin_program_service) queryWxpluginProgram(program_uuid string) (*tWxPluginProgram, error) {
	args := []interface{}{}
	args = append(args, program_uuid)
	tmp := tWxPluginProgram{}
	queryReq := &SqlQueryRequest{
		SQL:         "select program_id, program_uuid, program_name, appid, appsecrete, program_status, saler_uuid, program_type from t_wxplugin_program where program_uuid = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query wxpluginProgramByUuid[%s] error:%s", program_uuid, queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	return queryRep.Rows[0].(*tWxPluginProgram), nil
}
