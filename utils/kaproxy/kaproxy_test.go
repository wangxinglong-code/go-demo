package kaproxy

import (
	"encoding/json"
	"fmt"
	"go-demo/utils/common"
	"go-demo/utils/config"
	"go-demo/utils/logger"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	//初始化配置
	var err error
	configPath := "../../conf/config.toml"

	err = config.InitConfig(configPath)
	if err != nil {
	}

	//初始化日志
	err = logger.InitLoggerConfig()
	if err != nil {
	}
}

//go test -v -run TestKaProxyClient_Produce kaproxy_test.go
func TestKaProxyClient_Produce(t *testing.T) {
	var c *gin.Context
	c = new(gin.Context)
	c.Set("req_id", common.UniqueId())
	kaClient := NewClient(c, config.Config.KaProxy.Scheme, config.Config.KaProxy.Host, config.Config.KaProxy.Port, config.Config.KaProxy.Token)

	payload := map[string]string{
		"hello": "world",
	}
	payloadMsg, _ := json.Marshal(payload)
	msg := Message{Key: "test", Value: string(payloadMsg)}
	res, err := kaClient.ProduceWithHash(common.KaProxyDefaultTopic, msg)
	fmt.Println(err)
	fmt.Println(res)
}

//go test -v -run TestKaProxyClient_Consume kaproxy_test.go
func TestKaProxyClient_Consume(t *testing.T) {
	var c *gin.Context
	c = new(gin.Context)
	c.Set("req_id", common.UniqueId())

	kaClient := NewClient(c, config.Config.KaProxy.Scheme, config.Config.KaProxy.Host, config.Config.KaProxy.Port, config.Config.KaProxy.Token)
	res, err := kaClient.Consume(common.KaProxyDefaultGroup, common.KaProxyDefaultTopic, 0, 0)
	fmt.Println(err)
	fmt.Println(res)
}
