# go-demo
### 本地数据库依赖部署，推荐docker
#### mysql:
    docker pull mysql:5.7
    docker run -p 3306:3306 --name mysql  -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
#### redis:
    docker pull redis
    docker run -itd --name redis-test -p 6379:6379 redis 

#### 目录结构描述
```
├── Readme.md                   // 说明
├── .gitlab-ci.yml              // gitlab-ci
├── bin                         // build 二进制文件目录
├── conf                        // 配置文件目录
│   └── config.toml             // 配置toml文件, 可根据环境新增不同文件
├── controller                  // 项目控制器方法
│   └── api                     // 自定义控制器二级目录
│       └── apiTest.go          // 具体控制器逻辑代码
│   └── xxx                     // 可新增其他目录
├── docs                        // 项目文档，可编写接口文档和sql 文件等
├── middleware                  // 中间件目录
├── models                      // 数据库model文件目录
├── routers                     // 项目路由文件目录
├── services                    // 远程请求service目录
├── utils                       // 项目lib 库，包含配置映射、http 封装、数据库连接池等通用库
│   └── common                  // 通用配置目录。如通用函数、变量、方法等
│   └── config                  // 配置toml 映射
│   └── http                    // 通用http封装方法，如构造post,get等 
│   └── logger                  // 通用日志组件封装
│   └── mysql                   // 通用mysql 连接池封装
│   └── pgsql                   // 通用pgsql 连接池封装
│   └── redis                   // 通用redis 连接池封装
│   └── remote                  // 远程服务host/api目录
├── go.mod                      // go module 文件
├── main.go                     // 项目主入口文件，包含初始化、连接池等
├── Makefile                    // 可自定义编写脚本
└── Dockerfile                  // 项目Dockerfile文件
```

#### 项目启动
    go run main.go
    或者
    $ go build 
    $ ./bin/go-demo 
#### 指定具体配置文件
    go run main.go -f conf/config.toml

#### api访问方式
    ip:port
    eg: http://127.0.0.1:3001/v1/apiTest
    
#### docker build and run 
    $  make   //构建二进制包(注意不同平台build)
    $  docker build -t go-demo:test .   //build go-demo test 版本镜像
    $  docker run -p 3002:3001 --rm -it go-demo:test  //容器3001映射宿主机3002 
    http 访问127.0.0.1:3002