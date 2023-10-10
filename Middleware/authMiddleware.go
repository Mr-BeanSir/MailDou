package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"yxd/modle"
)

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	if !modle.IsOk(session) {
		session.Clear()
		session.Save()
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "用户登录失效，请重新登录",
			"tag":   "010",
		})
		return
	}
}
