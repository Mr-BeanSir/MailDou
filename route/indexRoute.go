package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexRoute(r *gin.Engine) {
	g := r.Group("/")
	{
		g.GET("/",
			func(c *gin.Context) {
				c.HTML(http.StatusOK, "index/index.html", gin.H{
					"people": "0+",
				})
			})
	}
}
