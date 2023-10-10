package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"yxd/controller"
)

func UserRoute(r *gin.Engine) {

	g := r.Group("/user")
	{
		// 登录/注册操作开始
		g.GET("/login", func(c *gin.Context) {
			session := sessions.Default(c)
			// 如果登录成功直接跳转到控制面板
			if session.Get("logged") == true {
				c.Header("location", "/dashboard")
				c.AbortWithStatus(http.StatusFound)
				return
			}
			var user controller.UserInfo
			var paramExit bool
			user.Username, paramExit = c.GetQuery("b")
			if !paramExit {
				c.HTML(http.StatusBadRequest, "index/error.html", gin.H{
					"error": "参数错误",
					"tag":   "032",
				})
				return
			}
			user.Id, paramExit = c.GetQuery("e")
			if !paramExit {
				c.HTML(http.StatusBadRequest, "index/error.html", gin.H{
					"error": "参数错误",
					"tag":   "033",
				})
				return
			}
			user.Time, paramExit = c.GetQuery("a")
			if !paramExit {
				c.HTML(http.StatusBadRequest, "index/error.html", gin.H{
					"error": "参数错误",
					"tag":   "034",
				})
				return
			}
			user.Token, paramExit = c.GetQuery("c")
			if !paramExit {
				c.HTML(http.StatusBadRequest, "index/error.html", gin.H{
					"error": "参数错误",
					"tag":   "035",
				})
				return
			}
			user.Regtoken = c.DefaultQuery("d", "nil")
			//  POST方法自动绑定参数使用
			//err := c.ShouldBind(&user)
			//if err != nil {
			//	c.HTML(http.StatusBadRequest, "index/error.html", gin.H{
			//		"error": "我还没见过这个报错",
			//		"tag":   "err",
			//	})
			//	return
			//}
			if bo, lUser := controller.SelectPrepare(user.Token, user.Id, user.Username); bo {
				//if !userLoginIsEmpty(user) {
				//	c.HTML(http.StatusOK, "index/error.html", gin.H{
				//		"error": "错误的传入参数",
				//		"tag":   "err",
				//	})
				//	return
				//}
				//logged := user.Login()
				//if logged["type"] == false {
				//	c.HTML(http.StatusOK, "index/error.html", gin.H{
				//		"error": logged["error"],
				//		"tag":   logged["tag"],
				//	})
				//	return
				//}
				//setUser, ok := logged["user"].(controller.Users)
				//if !ok {
				//	c.HTML(http.StatusOK, "index/error.html", gin.H{
				//		"error": "类型断言问题，类型转化失败",
				//		"tag":   "006",
				//	})
				//	return
				//}
				session.Set("id", lUser.Id)
				session.Set("username", lUser.Username)
				session.Set("token", lUser.Token)
				session.Set("logged", true)
				session.Save()
				c.Header("location", "/dashboard")
				c.AbortWithStatus(http.StatusFound)
				return
			} else {
				regBack := user.Reg()
				if !reflect.ValueOf(regBack["type"]).Bool() {
					c.HTML(http.StatusOK, "index/error.html", gin.H{
						"error": reflect.ValueOf(regBack["error"]).String(),
						"tag":   regBack["tag"],
					})
					return
				}
				setUser, ok := regBack["user"].(controller.UserInfo)
				if !ok {
					c.HTML(http.StatusOK, "index/error.html", gin.H{
						"error": "类型断言问题，类型转化失败",
						"tag":   "006",
					})
					return
				}
				session.Set("id", setUser.Id)
				session.Set("username", setUser.Username)
				session.Set("token", setUser.Token)
				session.Set("logged", true)
				session.Save()
				c.Header("location", "/dashboard")
				c.AbortWithStatus(http.StatusFound)
				return
			}
			c.HTML(http.StatusOK, "index/error.html", gin.H{
				"error": "请明确要进行的动作",
				"tag":   "007",
			})
			return
		}) // 登录/登录操作结束
		g.GET("/logout", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Clear()
			session.Save()
			c.Header("location", "/")
			c.AbortWithStatus(http.StatusFound)
			return
		}) // 登出操作结束
	}
}

func userLoginIsEmpty(u controller.UserInfo) bool {
	if reflect.DeepEqual(u.Id, "") {
		return false
	}
	if reflect.DeepEqual(u.Token, "") {
		return false
	}
	if reflect.DeepEqual(u.Username, "") {
		return false
	}
	return true
}

func userRegIsEmpty(u controller.UserInfo) bool {
	if reflect.DeepEqual(u.Id, "") {
		return false
	}
	if reflect.DeepEqual(u.Token, "") {
		return false
	}
	if reflect.DeepEqual(u.Time, "") {
		return false
	}
	if reflect.DeepEqual(u.Regtoken, "") {
		return false
	}
	if reflect.DeepEqual(u.Username, "") {
		return false
	}
	return true
}
