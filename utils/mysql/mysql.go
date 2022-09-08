package mysql

// https://gorm.io/zh_CN/docs/

import (
	"errors"
	"fmt"
	"go-demo/conf/beta"
	"go-demo/conf/local"
	"go-demo/conf/pre"
	"go-demo/conf/release"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"go-demo/utils/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//mysql 主从读写分离
var MysqlConnPool map[string][]*gorm.DB

var MysqlConf map[string]config.MysqlModuleConfig

//初始化mysql 入口
func InitMySQLPool() error {
	switch config.Config.Env {
	case "local":
		MysqlConf = local.MysqlConf
	case "pre":
		MysqlConf = pre.MysqlConf
	case "beta":
		MysqlConf = beta.MysqlConf
	case "release":
		MysqlConf = release.MysqlConf
	default:
		MysqlConf = pre.MysqlConf
	}

	MysqlConnPool = make(map[string][]*gorm.DB)
	for k, v := range MysqlConf {
		for i := 0; i < len(v.Nodes); i++ {
			MysqlConnPoolItem, err := InitMysqlConn(v.Nodes[i], v.Rule)
			MysqlConnPool[k] = append(MysqlConnPool[k], MysqlConnPoolItem)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//获取连接
func GetConn(aliasDb string, dbNode int, useMaster bool) (*gorm.DB, error) {
	if _, ok := MysqlConnPool[aliasDb]; !ok {
		return nil, errors.New("找不到mysql库别名")
	}
	if dbNode >= len(MysqlConnPool[aliasDb]) {
		return nil, errors.New("mysql选择节点超限")
	}

	if useMaster {
		return MysqlConnPool[aliasDb][dbNode].Clauses(dbresolver.Write), nil
	} else {
		return MysqlConnPool[aliasDb][dbNode], nil
	}
}

//初始化连接池
func InitMysqlConn(mysqlClusterConfig config.MysqlClusterConfig, mysqlRuleConfig config.MysqlRuleConfig) (*gorm.DB, error) {
	//Silent LogLevel = iota + 1 Error Warn Info
	var loggerLevel logger.LogLevel
	if config.GetMode() == "release" {
		loggerLevel = 1
	} else {
		loggerLevel = 4
	}
	mysqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      loggerLevel, // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)

	dsnMaster := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlClusterConfig.Master.User, mysqlClusterConfig.Master.Password, mysqlClusterConfig.Master.Host, mysqlClusterConfig.Master.Database)
	dsnSlave := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlClusterConfig.Slave.User, mysqlClusterConfig.Slave.Password, mysqlClusterConfig.Slave.Host, mysqlClusterConfig.Slave.Database)
	DB, err := gorm.Open(mysql.Open(dsnMaster), &gorm.Config{Logger: mysqlLogger})
	if err != nil {
		return nil, err
	}

	DB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(dsnSlave)},
	}).Register(dbresolver.Config{}).
		SetConnMaxIdleTime(mysqlRuleConfig.ConnMaxIdleTime).
		SetConnMaxLifetime(mysqlRuleConfig.ConnMaxLifetime).
		SetMaxIdleConns(mysqlRuleConfig.MaxIdleConns).
		SetMaxOpenConns(mysqlRuleConfig.MaxOpenConns))
	return DB, nil
}
