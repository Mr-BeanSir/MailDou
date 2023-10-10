## 邮小豆

![程序截图](https://i.mji.rip/2023/10/10/91bfc2a8c84dcf91a9e464ebec7ed599.png)

#### 如何使用？

```
git clone https://github.com/Mr-BeanSir/MailDou
cd MailDou
go mod tidy
```

新建数据库，导入根目录下maildou.sql

安装宝塔邮局插件，配置api和白名单

修改conf目录下 btApi.go、conf.go、sqlConn.go 文件

```
go build mail.go
mail.exe
```

#### 登录/注册

http://localhost:8080/user/login?a=时间戳&b=用户名&c=TOKEN&d=计算md5&e=用户ID

计算md5：conf-conf.go-regkey + TOKEN + 时间戳

## 说明

程序本意提供go语言有关邮箱的库使用及代码参考，不推荐公开搭建使用