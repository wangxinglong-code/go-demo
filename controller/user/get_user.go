package user

import (
	"go-demo/services/user"
	"go-demo/utils/common"
	"go-demo/utils/http"
	"go-demo/utils/logger"

	"github.com/gin-gonic/gin"
)

type GetUserParam struct {
	Id int `json:"id"`
}

func GetUser(c *gin.Context) {
	requestParam := &GetUserParam{}
	err := http.GetBodyParam(c, requestParam)
	if err != nil {
		logger.Debugf(c, "request param err:%+v, param:%+v", err.Error(), requestParam)
		http.ResponseErrorCodeAndData(c, common.ERR_INPUT_FMT, common.JsonEmptyObj)
		return
	}
	if requestParam.Id == 0 {
		http.ResponseErrorCodeAndData(c, common.ERR_INPUT, common.JsonEmptyObj)
		return
	}

	info, err := user.GetUser(requestParam.Id)
	if err != nil {
		logger.Infof(c, "get user error:%+v , param:%+v", err, requestParam.Id)
		http.ResponseErrorCodeAndData(c, common.ERR_MYSQL_GET_DATA, common.JsonEmptyObj)
		return
	}

	http.ResponseSuccess(c, info)
	return
}
