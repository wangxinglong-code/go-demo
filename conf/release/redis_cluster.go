package release

import (
	"go-demo/utils/config"
)

var RedisConf = map[string]config.RedisModuleConfig{
	"session": {
		Rule: config.RedisRuleConfig{
			MaxIdleConns:   200,
			IdleTimeOutSec: 240,
		},
		Nodes: map[int]config.RedisClusterConfig{
			0: {
				Master: []config.RedisDbConfig{
					0: {
						Host:     "127.0.0.1:6379",
						Password: "123456",
						SelectDb: 1,
					},
				},
				Slave: []config.RedisDbConfig{
					0: {
						Host:     "127.0.0.1:6379",
						Password: "123456",
						SelectDb: 1,
					},
				},
			},
		},
	},
}
