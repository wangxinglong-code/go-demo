package middleware

import (
	"go-demo/utils/common"
	"go-demo/utils/config"

	"github.com/gin-gonic/gin"
)

// 设置程序开始时间
func SetRequestStartTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("req_id") == "" {
			//生成请求id
			c.Set("req_id", common.UniqueId())
		} else {
			c.Set("req_id", c.Request.Header.Get("req_id"))
		}
		if c.Request.Header.Get("req_source") == "" {
			c.Set("req_source", config.Config.AppName)
		} else {
			c.Set("req_source", c.Request.Header.Get("req_source"))
		}
		c.Set("requestStartTime", common.Start())
		c.Next()
	}
}

//设置其他中间件
func TestFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Println("filter test")
		c.Next()
	}
}
