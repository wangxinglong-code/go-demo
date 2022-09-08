package routers

import (
	"fmt"
	"go-demo/controller/user"
	"go-demo/middleware"

	"github.com/gin-gonic/gin"
)

func Init(g *gin.Engine) (err error) {
	if g == nil {
		err = fmt.Errorf("nil gin engine")
		return err
	}

	//通用设置
	g.Use(middleware.SetRequestStartTime())
	//健康检查-k8s
	g.GET("/devops/status")

	//具体的路由
	baseGroup := g.Group("", middleware.TestFilter())
	{

		goDemoGroup := baseGroup.Group("/v1/go-demo/")
		{
			goDemoGroup.POST("demo", user.GetUser)
		}
	}

	return nil
}
