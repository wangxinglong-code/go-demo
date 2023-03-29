package pre

import (
	"go-demo/utils/config"
	"time"
)

var MysqlConf = map[string]config.MysqlModuleConfig{
	"go-demo": {
		Rule: config.MysqlRuleConfig{
			TableNums:       1,
			ConnMaxIdleTime: time.Second * 90, //最大空闲时长
			ConnMaxLifetime: time.Hour * 7,    //最大存活时长show variables like 'wait_timeout'
			MaxIdleConns:    100,
			MaxOpenConns:    200, //xshow variables like 'max_connections'
		},
		Nodes: map[int]config.MysqlClusterConfig{
			0: {
				Master: config.DbConfig{
					Host:     "127.0.0.1",
					Port:     3306,
					User:     "root",
					Password: "123456",
					Database: "go-demo-db",
				},
				Slave: config.DbConfig{
					Host:     "127.0.0.1",
					Port:     3306,
					User:     "root",
					Password: "123456",
					Database: "go-demo-db",
				},
			},
		},
	},
}
