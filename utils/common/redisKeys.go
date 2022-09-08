package common

//common config redis key prefix
const (
	COOKIE_TK              = "tk"
	COMMON_EXPIRE_TIME     = 259200 //60 * 60 * 24 * 3
	COMMON_TMP_EXPIRE_TIME = 100    //60 * 30
	COMMON_DAY_EXPIRE_TIME = 86400  //一天

	BASE = "base"

	// redisCount
	REDIS_COUNT = "counter"
)
