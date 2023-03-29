package pre

import (
	"go-demo/utils/config"
	"time"
)

var PgsqlConf = map[string]config.PgsqlModuleConfig{
	"go-demo": {
		Rule: config.PgsqlRuleConfig{
			TableNums:       1,
			ConnMaxIdleTime: time.Second * 90, //max idle time
			ConnMaxLifetime: time.Hour * 7,    //show variables like 'wait_timeout'
			MaxIdleConns:    100,
			MaxOpenConns:    200, //show variables like 'max_connections'
		},
		Nodes: map[int]config.PgsqlClusterConfig{
			0: {
				Master: config.DbConfig{
					Host:     "127.0.0.1",
					Port:     5432,
					User:     "go_demo",
					Password: "123456",
					Database: "go-demo-db",
				},
				Slave: config.DbConfig{
					Host:     "127.0.0.1",
					Port:     5432,
					User:     "go_demo",
					Password: "123456",
					Database: "go-demo-db",
				},
			},
		},
	},
}
