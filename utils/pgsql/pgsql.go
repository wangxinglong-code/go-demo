package pgsql

import (
	"errors"
	"fmt"
	"go-demo/conf/beta"
	"go-demo/conf/local"
	"go-demo/conf/pre"
	"go-demo/conf/release"
	"go-demo/utils/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var PgsqlConnPool map[string][]*gorm.DB

var PgsqlConf map[string]config.PgsqlModuleConfig

//初始化pgsql 入口
func InitMySQLPool() error {
	switch config.Config.Env {
	case "local":
		PgsqlConf = local.PgsqlConf
	case "pre":
		PgsqlConf = pre.PgsqlConf
	case "beta":
		PgsqlConf = beta.PgsqlConf
	case "release":
		PgsqlConf = release.PgsqlConf
	default:
		PgsqlConf = pre.PgsqlConf
	}

	PgsqlConnPool = make(map[string][]*gorm.DB)
	for k, v := range PgsqlConf {
		for i := 0; i < len(v.Nodes); i++ {
			PgsqlConnPoolItem, err := InitPgsqlConn(v.Nodes[i], v.Rule)
			PgsqlConnPool[k] = append(PgsqlConnPool[k], PgsqlConnPoolItem)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//获取连接
func GetConn(aliasDb string, dbNode int, useMaster bool) (*gorm.DB, error) {
	if _, ok := PgsqlConnPool[aliasDb]; !ok {
		return nil, errors.New("找不到pgsql库别名")
	}
	if dbNode >= len(PgsqlConnPool[aliasDb]) {
		return nil, errors.New("pgsql选择节点超限")
	}

	if useMaster {
		return PgsqlConnPool[aliasDb][dbNode].Clauses(dbresolver.Write), nil
	} else {
		return PgsqlConnPool[aliasDb][dbNode], nil
	}
}

//初始化连接池
func InitPgsqlConn(pgsqlClusterConfig config.PgsqlClusterConfig, pgsqlRuleConfig config.PgsqlRuleConfig) (*gorm.DB, error) {
	//Silent LogLevel = iota + 1 Error Warn Info
	var loggerLevel logger.LogLevel
	if config.GetMode() == "release" {
		loggerLevel = 1
	} else {
		loggerLevel = 4
	}
	pgsqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      loggerLevel, // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)

	dsnMaster := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", pgsqlClusterConfig.Master.User, pgsqlClusterConfig.Master.Password, pgsqlClusterConfig.Master.Host, pgsqlClusterConfig.Master.Database)
	dsnSlave := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", pgsqlClusterConfig.Slave.User, pgsqlClusterConfig.Slave.Password, pgsqlClusterConfig.Slave.Host, pgsqlClusterConfig.Slave.Database)
	DB, err := gorm.Open(postgres.Open(dsnMaster), &gorm.Config{Logger: pgsqlLogger})
	if err != nil {
		return nil, err
	}

	DB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(dsnSlave)},
	}).Register(dbresolver.Config{}).
		SetConnMaxIdleTime(pgsqlRuleConfig.ConnMaxIdleTime).
		SetConnMaxLifetime(pgsqlRuleConfig.ConnMaxLifetime).
		SetMaxIdleConns(pgsqlRuleConfig.MaxIdleConns).
		SetMaxOpenConns(pgsqlRuleConfig.MaxOpenConns))
	return DB, nil
}
