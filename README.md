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
│   ├── const.go
│   ├── const_test.go
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
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
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
│       ├── conf
│       │   ├── conf.go
│       │   └── conf_test.go
│       ├── cos
│       │   ├── client.go
│       │   ├── client_test.go
│       │   ├── gg.jpeg
│       │   └── mm.jpeg
│       ├── json_util.go
│       └── sms
│           ├── sms.go
│           └── sms_test.go
├── main.go
└── server
    ├── conf
    │   ├── config.go
    │   └── config_test.go
    ├── controller    // 也叫resource 层
    │   └── user_controller.go
    ├── dao     // 数据库连接池
    │   ├── dao.go
    │   ├── redis.go
    │   ├── redis_test.go
    │   ├── sql.go
    │   ├── sql_test.go
    │   └── sqlx.go
    ├── dto   // Data transfer object,  接收的参数获， 或者返回的结构体都可以在这里定义。
    │   └── common.go
    ├── middleware   /// 中间件
    │   ├── logger.go
    │   ├── middleware.go
    │   └── response.go
    ├── model          // 模型层， data access object(DAO) 定义在这里
    │   └── user_model.go
    ├── mq
    │   └── kafka
    │       ├── kafka.go
    │       └── kafka_test.go
    ├── router  // 路由层
    │   └── router.go
    ├── service   // 业务逻辑层， 也叫logic, operation 层
    │   └── user_service.go
    └── util
        ├── const.go
        ├── log.go
        └── log_test.go


```

- doc: api文档， 用[go-swagger](https://goswagger.io/tutorial/todo-list.html)写文档
    - [自动生成文档](https://juejin.im/post/5b05138cf265da0ba7701a37)

## 日志：
- 使用[logrus](https://github.com/sohlich/elogrus)记录日志
- 尽量使用`WithFields`打印日志, Feild中的字段要和当前打印日志的上下文相关， 如订单相关的日志如下打印

```go
Log.WithFields(logrus.Fields{"order_id": 123456, "user_id": 1}).Errorf("订单付款失败: err: %s", err)
```

## 单元测试：
 
- 所有单元测试要不能依赖其他包，需要单独可以运行（见阿里java开发手册）。 

## 配置和地址:

- mongo, mysql, redis的主机， 端口， 账户， 密码 见zookeeper 的`superconf/union`

# go web 项目模版
 - https://github.com/eddycjy/go-gin-example
 - https://github.com/Keegan-y/gin_scaffold#%E6%96%87%E4%BB%B6%E5%88%86%E5%B1%82
 - https://github.com/e421083458/gin_scaffold
 - https://github.com/e421083458/go_gateway (这个项目的路由，和controller分层很值得学习)
 
 
# 运行： 
```bash
go get -u github.com/swaggo/swag/cmd/swag 
swag init

```