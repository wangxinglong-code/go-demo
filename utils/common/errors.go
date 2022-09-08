package common

import "errors"

//错误类
type Error struct {
	ErrorCode int
	ErrorMsg  string
}

func (this *Error) Error() string {
	return this.ErrorMsg
}

var (
	ERR_SUC         = &Error{ErrorCode: 0, ErrorMsg: "OK"}
	ERR_INPUT       = &Error{ErrorCode: 1001, ErrorMsg: "缺少参数"}
	ERR_INPUT_FMT   = &Error{ErrorCode: 1002, ErrorMsg: "参数格式错误"}
	ERR_INPUTANDMAX = &Error{ErrorCode: 1003, ErrorMsg: "缺少参数或参数超过上限"}

	//配置or链接
	ERR_MYSQL           = &Error{ErrorCode: 1201, ErrorMsg: "MySQL数据库报错"}
	ERR_MONGODB         = &Error{ErrorCode: 1202, ErrorMsg: "MongoDB数据库报错"}
	ERR_REDIS           = &Error{ErrorCode: 1203, ErrorMsg: "Redis数据库报错"}
	ERR_UNMARSHAL_ERROR = &Error{ErrorCode: 1204, ErrorMsg: "解析数据异常"}
	ERR_MARSHAL_ERROR   = &Error{ErrorCode: 1205, ErrorMsg: "解析数据异常"}
	ERR_GRPC            = &Error{ErrorCode: 1206, ErrorMsg: "Grpc连接报错"}

	//请求
	ERR_REMOTE_CURL        = &Error{ErrorCode: 1401, ErrorMsg: "远程请求错误"}
	ERR_REQUEST_OVER_LIMIT = &Error{ErrorCode: 1403, ErrorMsg: "请求超限"}

	//mysql
	ERR_MYSQL_CREATE_DATA       = &Error{ErrorCode: 1501, ErrorMsg: "创建数据失败"}
	ERR_MYSQL_CREATE_DATA_EXIST = &Error{ErrorCode: 15011, ErrorMsg: "创建/更新数据失败,数据已存在"}
	ERR_MYSQL_UPDATE_DATA       = &Error{ErrorCode: 1502, ErrorMsg: "更新数据失败"}
	ERR_MYSQL_DEL_DATA          = &Error{ErrorCode: 1503, ErrorMsg: "删除数据失败"}
	ERR_MYSQL_DEL_DATA_NOT      = &Error{ErrorCode: 15031, ErrorMsg: "删除数据失败,数据已绑定"}
	ERR_MYSQL_GET_DATA          = &Error{ErrorCode: 1504, ErrorMsg: "获取数据失败"}
	ERR_MYSQL_SELECT_DATA       = &Error{ErrorCode: 1505, ErrorMsg: "获取数据为空"}
	ERR_MYSQL_NAME_IS_EXISTS    = &Error{ErrorCode: 1506, ErrorMsg: "名称已存在"}

	//redis
	ERR_REDIS_CREATE_DATA = &Error{ErrorCode: 1601, ErrorMsg: "创建数据失败"}
	ERR_REDIS_UPDATE_DATA = &Error{ErrorCode: 1602, ErrorMsg: "更新数据失败"}
	ERR_REDIS_DEL_DATA    = &Error{ErrorCode: 1603, ErrorMsg: "删除数据失败"}
	ERR_REDIS_GET_DATA    = &Error{ErrorCode: 1604, ErrorMsg: "获取数据失败"}

	//grpc
	ERR_GRPC_SELECT_DATA = &Error{ErrorCode: 1701, ErrorMsg: "grpc获取数据失败"}
)

var (
	DB_NO_ROWS_AFFECTED = errors.New("no rows affected") // 没有数据修改错误
)
