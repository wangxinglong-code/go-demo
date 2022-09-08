package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

var Config Configure

type Configure struct {
	Port       int
	PprofPort  int
	RpcPort    int
	Env        string
	MonitorUrl string
	APPURL     string
	AppName    string
	Mode       string
	WarningKey string
	Log        LogConf
	KaProxy    KaProxyConfig
	RemoteHost RemoteHostConfig
}

type RemoteHostConfig struct {
	ServiceWebsocket            string
	ServiceEncyclopediaHTMLHost string
}

type KaProxyConfig struct {
	Scheme string
	Host   string
	Port   int
	Token  string
}

type LogConf struct {
	LogLevel           string
	LogFileStdout      string
	LogFileErrorStdout string
	AppKey             string
}

type DbConfig struct {
	Host     string
	Password string
	User     string
	Database string
}

/************    mysql配置结构体    ************/
type MysqlModuleConfig struct {
	Rule  MysqlRuleConfig
	Nodes map[int]MysqlClusterConfig
}

type MysqlRuleConfig struct {
	TableNums       uint32
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

type MysqlClusterConfig struct {
	Master DbConfig
	Slave  DbConfig
}

/************    pgsql配置结构体    ************/
type PgsqlModuleConfig struct {
	Rule  PgsqlRuleConfig
	Nodes map[int]PgsqlClusterConfig
}

type PgsqlRuleConfig struct {
	TableNums       uint32
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

type PgsqlClusterConfig struct {
	Master DbConfig
	Slave  DbConfig
}

/************    redis配置结构体    ************/
type RedisModuleConfig struct {
	Rule  RedisRuleConfig
	Nodes map[int]RedisClusterConfig
}

type RedisRuleConfig struct {
	MaxIdleConns   int
	IdleTimeOutSec int
}

//多主多从
type RedisClusterConfig struct {
	Master []RedisDbConfig
	Slave  []RedisDbConfig
}

type RedisDbConfig struct {
	Host     string
	Password string
	SelectDb int
}

func InitConf() error {
	var err error
	configPath := flag.String("f", "./conf/config.toml", "conf file")
	flag.Parse()

	//初始化配置文件
	err = InitConfig(*configPath)
	if err != nil {
		return err
	}
	return nil
}

func InitConfig(configPath string) error {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	if _, err := toml.Decode(string(configBytes), &Config); err != nil {
		return err
	}

	return nil
}

func GetMode() string {
	return Config.Mode
}

func GetPort() string {
	return fmt.Sprintf(":%d", Config.Port)
}
