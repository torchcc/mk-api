# mk-api
微信公众号

## 编程规范：参考阿里java开发手册，一下挑几点和以往习惯不同的

- 数据库名用名词单数！
- 包名用名词单数！



## ORM
- 可以使用GORM， 也可以不用，看自己喜好， 可以混着用。


## 项目结构

```
.
├── FAQ.md
├── README.md
├── deployment  
│   ├── local
│   │   ├── Dockerfile
│   │   ├── deploy.json
│   │   └── superconf.json
│   ├── prod
│   │   ├── Dockerfile
│   │   ├── deploy.json
│   │   └── superconf.json
│   └── test
│       ├── Dockerfile
│       ├── deploy.json
│       └── superconf.json
├── go.mod
├── go.sum
├── library
│   ├── README.md
│   ├── ecode
│   │   ├── common_ecode.go
│   │   ├── ecode.go
│   │   └── ecode_test.go
│   ├── superconf
│   │   └── config.go
│   └── util
│       └── json_util.go
├── main.go
└── server
    ├── conf
    │   ├── config.go
    │   └── config_test.go
    ├── controller
    │   ├── admin.go
    │   ├── admin_login.go
    │   └── user.go
    ├── dao
    │   ├── dao.go
    │   ├── redis.go
    │   ├── redis_test.go
    │   ├── sql.go
    │   └── sql_test.go
    ├── doc
    ├── log.go
    ├── log_test.go
    ├── logic
    ├── middleware
    │   └── middleware.go
    ├── model
    ├── mq
    │   └── kafka
    │       ├── kafka.go
    │       └── kafka_test.go
    ├── router
    │   └── router.go
    ├── util
    └── verify
```

- library: 公共库, 可以跨服务用的放在里面， 如全局错误码定义， 可以跨服务用的工具

- server: api服务主目录
- conf： 配置的注册
- controller 也就是resource 层
- dao Database Access Object， 创建redis， mysql， mongo的全局实例, 注意和model层区分
    - dao中mysql，redis的使用见对应的test文件。 
- doc: api文档， 用[go-swagger](https://goswagger.io/tutorial/todo-list.html)写文档
    - [自动生成文档](https://juejin.im/post/5b05138cf265da0ba7701a37)
- logic 也叫做operation层 
- middleware: 中间件， 如authentication（login_requied）, permission，csrf 等。 
- model 模型层
    - DO（数据对象）的定义， sql 语句， 与数据库的交互在着一层完成
    - 这层不打印错误日志， 错误日志在 logic层打印， 并且在controller层之前要捕捉处理完
- mq： kafka生产者和消费者client的封装。 
- route: 路由层
- util: 本项目的静态常量，本项目的工具类， 工具函数
- verify : 参数校验层， 可选(看看有没有必要增加着一层)


## 日志：
- 使用[logrus](https://github.com/sohlich/elogrus)记录日志
- 尽量使用`WithFields`打印日志, Feild中的字段要和当前打印日志的上下文相关， 如订单相关的日志如下打印

```go
Log.WithFields(logrus.Fields{"order_id": 123456, "user_id": 1}).Errorf("订单付款失败: err: %s", err)
```

## 单元测试：
 
- 所有单元测试要不能依赖其他包，需要单独可以运行（见阿里java开发手册）。 

## 配置和地址:

- [zookeeper](http://106.53.124.190:9090/login)
    - 账户： `admin`
    - 密码: `maikang`
    
- mongo, mysql, redis的主机， 端口， 账户， 密码 见zookeeper 的`superconf/union`


# go web 项目模版
 - https://github.com/eddycjy/go-gin-example
 - https://github.com/Keegan-y/gin_scaffold#%E6%96%87%E4%BB%B6%E5%88%86%E5%B1%82
 - https://github.com/e421083458/gin_scaffold
 - https://github.com/e421083458/go_gateway (这个项目的路由，和controller分层很值得学习)
 
 