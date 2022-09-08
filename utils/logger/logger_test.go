package logger

import (
	"fmt"
	"go-demo/utils/config"
	"go-demo/utils/mysql"
	"log"
	"os"
	"testing"
)

func init() {
	//初始化配置
	var err error
	configPath := "../../conf/config.toml"

	err = config.InitConfig(configPath)
	if err != nil {
	}

	err = mysql.InitMySQLPool()
	if err != nil {
	}
}

func TestWriteLog(t *testing.T) {
	pathTmp := "/www/arachnia_log/" + config.Config.Log.AppKey + "/bigdata"
	_, err := os.Stat(pathTmp)
	fmt.Println(err)
	if err != nil {
		err := os.MkdirAll(pathTmp, 0777)
		if err != nil {
			log.Fatal(err)
		}
		os.Chmod(pathTmp, 0777) //通过chmod重新赋权限
	}

	BigInfow("aaaa")
}
