package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminNewRouter(g *gin.RouterGroup) {
	g.GET("/news", controller.IndexNew)
	g.POST("/new", controller.StoreNew)
	g.GET("/new/:id", controller.ShowNew)
	g.PUT("/new/:id", controller.UpdateNew)
	g.DELETE("/new/:id", controller.DeleteNew)
	g.PUT("/new/:id/status", controller.ToggleNewStatus)
	g.PUT("/new/:id/topstatus", controller.ToggleNewTopStatus)
}
