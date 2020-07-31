package handler

import (
	"context"
	"encoding/json"
	"my-filestore-server/db"
	"my-filestore-server/service/file/proto"
	"my-filestore-server/util"
)

type File struct{}

func (f *File) FileQuery(ctx context.Context, req *proto.ReqFileQuery, resp *proto.RespFileQuery) error {
	username := req.Username
	limit := req.Limit

	userFiles, err := db.QueryUserFileMetas(username, int(limit))
	if err != nil {
		resp.Code = util.StatusQueryUserFilesError
		resp.Message = "查询 user file表失败"
		return nil
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		resp.Code = util.StatusQueryUserFilesError
		resp.Message = "数据格式有误"
		return nil
	}

	resp.Code = util.StatusOK
	resp.Message = "OK"
	resp.Data = data
	return nil
}
