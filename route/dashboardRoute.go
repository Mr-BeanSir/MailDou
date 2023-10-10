package route

import (
	"github.com/gin-gonic/gin"
	middleware "yxd/Middleware"
	"yxd/controller"
)

func DashboardRoute(r *gin.Engine) {
	g := r.Group("/dashboard", middleware.Auth)
	{
		g.GET("", controller.DashboardIndex)
		g.POST("/c_email", controller.DashboardCreateEmail)
		g.GET("/d_email/:email", controller.DashboardDeleteEmail)
		g.POST("transfer", controller.DashboardTransfer)
	}
}
