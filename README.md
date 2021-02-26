## 功能

User Story 1

- [x] 用户注册、登录
- [x] 发布问题
- [x] 查看问题列表

User Story 2

- [x] 修改/删除问题
- [x] 回答问题
- [x] 修改回答
- [x] 删除回答

User Story 3

- [x] 热门问题列表
- [x] 对回答点赞和踩
- [x] 个人中心查看自己已发布的问题
- [x] 个人中心查看自己的回答
- [x] 个人中心查看自己点赞的回答
- [x] 修改个人信息、头像

附加功能

- [x] 图片上传
- [x] 前端页面（https://github.com/recovic/im-zhihu-vue）
- [ ] 第三方登录
- [ ] 高并发

## 环境依赖

- [Gin](https://github.com/gin-gonic/gin): 轻量级Web框架
- [GORM](http://gorm.io/docs/index.html): ORM工具，连接操作数据库
- [Go-Redis](https://github.com/go-redis/redis): 连接 redis 服务器用作缓存
- [go-model](https://gopkg.in/jeevatkm/go-model.v1)：用于结构体间的自动映射
- [Jwt-Go](https://github.com/dgrijalva/jwt-go): JWT组件，用作身份验证
- [cron](https://github.com/robfig/cron): 定时任务库，用于缓存同步
- [logrus](https://github.com/sirupsen/logrus)：记录日志

## 目录结构

```
├── cache            redis 缓存相关
├── config           项目的静态配置
├── controller       API 到功能的映射
├── dto              与前端交互实体的定义
├── enum             枚举类
├── middleware       中间件
├── repository       数据库模型以及相关操作
├── result           结果类及错误信息的统一定义
├── service          将比较复杂的业务从 controller 层分离出来
├── tool             工具类
| main.go            项目入口
```

## 构建代码

go 版本：1.15

构建及下载依赖：
```
go build
```

运行：
```
./imitate-zhihu
```

服务器环境变量：
```
IZ_ENV_MODE=release
IZ_DB_ADDR=127.0.0.1:3306    //数据库地址
IZ_DB_USERNAME=xxx      //数据库用户名
IZ_DB_PASSWORD=xxx      //数据库密码
IZ_JWT_SECRET=abcd1234  //JWT Secret（任意字符串）
```

运行前先导入 sql 中的结构，导入到 zhihu 数据库用于运行，导入到 zhihu_test 数据库用于单元测试。