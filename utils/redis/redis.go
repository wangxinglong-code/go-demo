package redis

/**
 *
 */

import (
	"errors"
	"fmt"
	"go-demo/conf/beta"
	"go-demo/conf/local"
	"go-demo/conf/pre"
	"go-demo/conf/release"
	"go-demo/utils/config"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	PROTOCOL = "tcp" //connection protocol
)

//模块名-cluster_node-db_node
var RedisConnMasterPool map[string][][]*redis.Pool
var RedisConnSlavePool map[string][][]*redis.Pool
var RedisConf map[string]config.RedisModuleConfig

//初始化redis session
func InitRedisPool() error {
	switch config.Config.Env {
	case "local":
		RedisConf = local.RedisConf
	case "pre":
		RedisConf = pre.RedisConf
	case "beta":
		RedisConf = beta.RedisConf
	case "release":
		RedisConf = release.RedisConf
	default:
		RedisConf = pre.RedisConf
	}

	RedisConnMasterPool = make(map[string][][]*redis.Pool)
	RedisConnSlavePool = make(map[string][][]*redis.Pool)
	for k, v := range RedisConf {
		for i := 0; i < len(v.Nodes); i++ {
			nodes := v.Nodes[i]
			var RedisConnPoolMasterItem []*redis.Pool
			for j := 0; j < len(nodes.Master); j++ {
				RedisConnPoolItem := initRedisPool(nodes.Master[j].Host, nodes.Master[j].Password, v.Rule.MaxIdleConns, v.Rule.IdleTimeOutSec, nodes.Master[j].SelectDb)
				RedisConnPoolMasterItem = append(RedisConnPoolMasterItem, RedisConnPoolItem)
			}
			RedisConnMasterPool[k] = append(RedisConnMasterPool[k], RedisConnPoolMasterItem)

			var RedisConnPoolSlaveItem []*redis.Pool
			for j := 0; j < len(nodes.Slave); j++ {
				RedisConnPoolItem := initRedisPool(nodes.Slave[j].Host, nodes.Slave[j].Password, v.Rule.MaxIdleConns, v.Rule.IdleTimeOutSec, nodes.Slave[j].SelectDb)
				RedisConnPoolSlaveItem = append(RedisConnPoolSlaveItem, RedisConnPoolItem)
			}
			RedisConnSlavePool[k] = append(RedisConnSlavePool[k], RedisConnPoolSlaveItem)
		}
	}

	return nil
}

/**
 * Redis Pool
 * server 127.0.0.1:6379
 * IdleTimeout  超时
 * MaxIdle 连接池最大容量
 * MaxActive 最大活跃数量
 * dbno 选择db127.0.0.1:6379:password:1
 *
 */
func initRedisPool(server, password string, maxIdle, idleTimeout int, selectDb int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(PROTOCOL, server)
			if err != nil {
				return nil, err
			}

			if len(password) != 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if selectDb >= 0 {
				if _, err := c.Do("SELECT", selectDb); err != nil {
					c.Close()
					fmt.Println(err)
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if t.Add(time.Duration(maxIdle) * time.Second).After(time.Now()) {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

//获取主从节点，适用cluster db 下多个主从
func getDbNode(endNum int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(endNum)
}

//获取连接池
func GetConnPool(aliasName string, dbNode int, useMaster bool) (*redis.Pool, error) {
	if useMaster { //主库
		if _, ok := RedisConnMasterPool[aliasName]; !ok {
			return nil, errors.New("找不到redis库别名")
		}
		if dbNode >= len(RedisConnMasterPool[aliasName]) {
			return nil, errors.New("redis选择节点超限")
		}
		return RedisConnMasterPool[aliasName][dbNode][getDbNode(len(RedisConnMasterPool[aliasName][dbNode]))], nil
	} else {
		if _, ok := RedisConnSlavePool[aliasName]; !ok {
			return nil, errors.New("找不到redis库别名")
		}
		if dbNode >= len(RedisConnSlavePool[aliasName]) {
			return nil, errors.New("redis选择节点超限")
		}
		return RedisConnSlavePool[aliasName][dbNode][getDbNode(len(RedisConnSlavePool[aliasName][dbNode]))], nil
	}

}

func GetStringFromRedis(aliasName string, dbNode int, key string) (string, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return "", err
	}
	conn := pool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do("get", key))
	if err != nil {
		return "", err
	}

	return value, nil
}

func GetUint64FromRedis(aliasName string, dbNode int, key string) (uint64, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	value, err := redis.Uint64(conn.Do("get", key))
	if err != nil {
		return 0, err
	}

	return value, nil
}

func GetIntFromRedis(aliasName string, dbNode int, key string) (int, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("get", key))
	if err != nil {
		return 0, err
	}

	return value, nil
}

func SetStringToRedis(aliasName string, dbNode int, key, value string) error {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	_, err = redis.String(conn.Do("set", key, value))
	if err != nil {
		return err
	}

	return nil
}

func SetStringToRedisEX(aliasName string, dbNode int, key, value string, exTime int64) error {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	_, err = redis.String(conn.Do("set", key, value, "EX", exTime))
	if err != nil {
		return err
	}

	return nil
}

func DelKeyFromRedis(aliasName string, dbNode int, key string) (int64, error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("del", key))
}

func SetNXWithExpireToRedis(aliasName string, dbNode int, key, value string, expire int) (string, error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return "", err
	}
	conn := pool.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("set", key, value, "EX", expire, "NX"))
	if err != nil {
		return r, err
	}

	return r, nil
}

func SetNXToRedis(aliasName string, dbNode int, key, value string) (string, error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return "", err
	}
	conn := pool.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("set", key, value, "NX"))
	if err != nil {
		return r, err
	}

	return r, nil
}

//递增
func Incr(aliasName string, dbNode int, key string) error {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("incr", key)
	return err
}

func Setex(aliasName string, dbNode int, key string, seconds int64, value string) error {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("setex", key, seconds, value)
	return err
}

func HsetStringToRedis(aliasName string, dbNode int, key, field, value string) (int, error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("hset", key, field, value))
}

func ExpireKeyToRedis(aliasName string, dbNode int, key string, expire int) (int, error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("expire", key, expire))
}

func HgetStringToRedis(aliasName string, dbNode int, key, field string) (string, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return "", err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("hget", key, field))
}

func HdelStringToRedis(aliasName string, dbNode int, key, field string) (int, error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("hdel", key, field))
}

func SAdd(aliasName string, dbNode int, key, field string) (int, error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("SADD", key, field))
}

func Scard(aliasName string, dbNode int, key string) (int, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("SCARD", key))
}

func Smembers(aliasName string, dbNode int, key string) ([]interface{}, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return nil, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Values(conn.Do("SMEMBERS", key))
}

func SmembersToStrings(aliasName string, dbNode int, key string) ([]string, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return nil, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("SMEMBERS", key))
}

func Sdiffstore(aliasName string, dbNode int, key, firstKey, twoKey string) (int, error) {
	pool, err := GetConnPool(aliasName, dbNode, false)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("SDIFFSTORE", key, firstKey, twoKey))
}

func Ttl(aliasName string, dbNode int, key string) (second int64, err error) {
	pool, err := GetConnPool(aliasName, dbNode, true)
	if err != nil {
		return 0, err
	}
	conn := pool.Get()
	defer conn.Close()
	second, err = redis.Int64(conn.Do("ttl", key))
	return second, err
}
