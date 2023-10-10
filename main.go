package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"yxd/route"
)

func main() {
	r := gin.Default()
	//设置session基础配置
	store := cookie.NewStore([]byte("test"))
	r.Use(sessions.Sessions("yxdou", store))
	//读入模板文件
	r.LoadHTMLGlob("./template/**/*")
	//设置静态文件路径索引
	r.Static("/static", "./static")
	r.Static("/index", "./index")
	//设置首页路由
	route.IndexRoute(r)
	//设置登录注册页路由
	route.UserRoute(r)
	//注册控制面板路由
	route.DashboardRoute(r)
	//注册邮件阅读相关路由
	route.Mail(r)
	r.Run(":8080")
}
