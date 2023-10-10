package route

import (
	"github.com/gin-gonic/gin"
	middleware "yxd/Middleware"
	"yxd/controller"
)

func Mail(r *gin.Engine) {
	g := r.Group("/mail", middleware.Auth)
	{
		g.POST("", controller.IndexMail)
		g.GET("/list/:type", controller.ListMail)
		g.GET("/read/:type/:id", controller.ReadMail)
		g.GET("/del/:type/:id", controller.DelMail)
		g.GET("/content/:type/:id", controller.Html)
	}
}
