package app

import (
	"context"
	rpcLogger "go-demo/grpc/logger"
	"go-demo/grpc/protobuf/app/go_demo"
	"go-demo/models"
	"go-demo/utils/common"
	"time"

	"gorm.io/gorm"
)

type Server struct{}

func (s *Server) GetUser(ctx context.Context, in *go_demo.GetUserRequest) (*go_demo.GetUserReply, error) {
	reqTime := time.Now()
	res := &go_demo.GetUserReply{}
	if in.Uid == 0 {
		return res, common.ERR_INPUT
	}

	user, err := new(models.User).Get(int(in.Uid))
	if err != nil && err != gorm.ErrRecordNotFound {
		return res, common.ERR_MYSQL_GET_DATA
	}
	if err == gorm.ErrRecordNotFound {
		return res, nil
	}

	res.Uid = int32(user.Uid)
	res.Name = user.Name
	res.Age = int32(user.Age)

	rpcLogger.RpcInfow(in.ReqSource, in.ReqId, "GetUser", reqTime, in, res)
	return res, nil
}
