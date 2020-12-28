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

运行前先导入 sql 中的结构，导入到 zhihu 用于运行，导入到 zhihu_test 用于单元测试。