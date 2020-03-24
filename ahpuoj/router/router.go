package router

import (
	"ahpuoj/controller"
	"ahpuoj/middleware"

	"github.com/gin-gonic/gin"
)

func handelRouterGroup() {

}

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/login", controller.Login)
	router.POST("/api/register", controller.Register)
	router.POST("/api/findpass", controller.SendFindPassEmail)
	router.GET("/api/verifyresetpasstoken", controller.VeriryResetPassToken)
	router.POST("/api/resetpassbytoken", controller.ResetPassByToken)

	// 添加解析token中间件
	router.Use(middleware.ParseTokenMiddleware())
	// 无需用户登录的api
	nologin := router.Group("/api")
	ApiNologinRouter(nologin)
	// 添加JWT认证中间件
	router.Use(middleware.JwtauthMiddleware())
	// 需要用户登录的api
	user := router.Group("/api")
	ApiUserRouter(user)
	// 后台路由组 添加Casbin权限控制中间件
	admin := router.Group("/api/admin", middleware.CasbinMiddleware())
	// 后台路由
	ApiAdminRouter(admin)

	return router
}
