package common

import "os"

//获取本机hostname
var HostName string = "go-demo"

func InitHostName() error {
	HostName, _ = os.Hostname()
	return nil
}

func GetHostNamePrefix() string {
	return HostName + "	"
}
