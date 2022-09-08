package main

import (
	"context"
	"fmt"
	rpcServer "go-demo/grpc/server"
	"go-demo/routers"
	"go-demo/utils/config"
	httpUtils "go-demo/utils/http"
	"go-demo/utils/logger"
	"go-demo/utils/mysql"
	"go-demo/utils/pgsql"
	"go-demo/utils/redis"

	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	err := Init()
	if err != nil {
		panic(err)
	}

	//设置模式
	gin.SetMode(config.GetMode())

	router := gin.Default()
	router.Use(httpUtils.Cors())

	//设置session
	////设置session--redis必须配置session模块库
	//redisPool, err := redis.GetConnPool("session", 0, true)
	//if err != nil {
	//	panic(err)
	//}
	//store, err := sessionRedis.NewStoreWithPool(redisPool, []byte("go-demo"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//store.Options(sessions.Options{
	//	MaxAge: int(time.Now().Unix()) + int(30*time.Minute), // 30min
	//	Path:   "/",
	//})
	//
	//router.Use(sessions.Sessions("go-demo_session", store))

	// 注册路由
	err = routers.Init(router)
	if err != nil {
		logger.Panicf("routers.Init %s error", err.Error())
		panic(err)
	}

	//pprof
	go func() {
		log.Println(http.ListenAndServe(":"+fmt.Sprint(config.Config.PprofPort), nil))
	}()

	port := config.GetPort()

	srv := &http.Server{
		Addr:           port,
		Handler:        router,
		IdleTimeout:    120 * time.Second,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.SystemInfof("Server Starting")

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panicf("listen port %s error", port)
		}
	}()

	gracefulExit(srv)
}

/**
 * 释放连接池资源，优雅退出
 * @param  {[type]} srv *http.Server  [description]
 * @return {[type]}     [description]
 */
func gracefulExit(srv *http.Server) {
	quit := make(chan os.Signal)
	// signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for s := range quit {
		switch s {
		// case syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2:

		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

			if err := srv.Shutdown(ctx); err != nil {
				logger.SystemInfof("Server Shutdown:", err)
			}
			logger.SystemInfof("Server Exiting")

			cancel()
			os.Exit(0)
		}
	}
}

/**
 * 初始化 日志，mysql，redis等
 */

func Init() error {
	var err error

	//初始化配置
	err = config.InitConf()
	if err != nil {
		return err
	}

	//初始化日志
	err = logger.InitLoggerConfig()
	if err != nil {
		return err
	}

	//初始化MySQL
	err = mysql.InitMySQLPool()
	if err != nil {
		return err
	}

	//初始化PgSQL
	err = pgsql.InitPgSQLPool()
	if err != nil {
		return err
	}

	//初始化Redis
	err = redis.InitRedisPool()
	if err != nil {
		return err
	}

	//初始化rpc连接池
	//err = rpcClient.InitRpcClient()
	//if err != nil {
	//	return err
	//}

	//初始化Grpc--按需，
	go func() {
		rpcServer.InitGrpc()
	}()
	return nil
}
