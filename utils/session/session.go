package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 设置session
func SetSessionInfo(c *gin.Context, key string, value interface{}) (err error) {
	session := sessions.Default(c)
	session.Set(key, value)
	return nil
}

func SaveSession(c *gin.Context) error {
	session := sessions.Default(c)
	return session.Save()
}

// 获取session
func GetSessionInfo(c *gin.Context, key string) interface{} {
	return sessions.Default(c).Get(key)
}

// flush
func FlushSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}
