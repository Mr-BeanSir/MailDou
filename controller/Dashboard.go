package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"yxd/Dao"
	"yxd/conf"
)

func DashboardIndex(c *gin.Context) {
	session := sessions.Default(c)
	emails, err := Dao.SelectAllEmail(session)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "获取邮箱失败",
			"tag":   "011",
		})
		return
	}
	c.HTML(http.StatusOK, "dashboard/index.html", gin.H{
		"emails":   emails,
		"username": session.Get("username"),
		"id":       session.Get("id"),
	})

}

func DashboardCreateEmail(c *gin.Context) {
	session := sessions.Default(c)
	prefix, err := c.GetPostForm("prefix")
	if !err {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "传递参数有误",
			"tag":   "015",
		})
		return
	}
	address := prefix + "@maildou.cc"
	_, err1 := Dao.CreateEmail(session, address, prefix)
	fmt.Println(err1)
	if err1 != nil {
		fmt.Println("here")
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": err1,
			"tag":   "01X",
		})
		return
	}
	c.Header("location", "/dashboard")
	c.AbortWithStatus(http.StatusFound)

}

func DashboardDeleteEmail(c *gin.Context) {
	session := sessions.Default(c)
	email := c.Param("email")
	if !Dao.Is_Self(session.Get("id").(string), email) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱不属于你，你tm想干坏事",
			"tag":   "110",
		})
		return
	}
	if Dao.DeleteEmail(session, email) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "删除邮箱成功",
			"tag":   "-1",
		})
		return
	}
}

func DashboardTransfer(c *gin.Context) {
	session := sessions.Default(c)
	email, err := c.GetPostForm("email")
	if !err {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "传递参数有误",
			"tag":   "036",
		})
		return
	}
	match, _ := regexp.Match("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", []byte(email))
	if !match {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱匹配错误",
			"tag":   "035",
		})
		return
	}
	pwd, err := c.GetPostForm("pwd")
	if !err {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "传递参数有误",
			"tag":   "015",
		})
		return
	}
	pid, err := c.GetPostForm("pid")
	if !err {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "传递参数有误",
			"tag":   "037",
		})
		return
	}
	if !transferIsEmpty(email, pwd, pid) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "传递参数有误",
			"tag":   "033",
		})
		return
	}
	if !Dao.Is_Self(session.Get("id").(string), email) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱不属于你，你tm想干坏事",
			"tag":   "110",
		})
		return
	}
	if !Dao.Is_Email(email, pwd) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱与密码错误",
			"tag":   "111",
		})
		return
	}
	b, transferUserName := Dao.Is_Pid(pid)
	if !b {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "要移动的用户不存在",
			"tag":   "112",
		})
		return
	}

	if !conf.ChangePassword(email) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "重置邮箱密码失败，转移失败",
			"tag":   "113",
		})
		return
	}

	if !Dao.Transfer(pid, email, transferUserName) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "转移邮箱所有人失败咯",
			"tag":   "115",
		})
		return
	}
	c.HTML(http.StatusOK, "index/true.html", gin.H{
		"error": "转移成功！请返回",
		"tag":   "-1",
	})
	return
}

func transferIsEmpty(email string, pwd string, pid string) bool {
	if email == "" || pwd == "" || pid == "" {
		return false
	}
	return true
}
