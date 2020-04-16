package middleware

import (
	"ahpuoj/model"
	"ahpuoj/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sub1 string
		var sub2 string
		user, _ := c.Get("user")
		if user, ok := user.(model.User); ok {
			sub1 = user.Role
			sub2 = strconv.Itoa(user.Id)
		}

		obj := c.Request.URL.Path
		act := c.Request.Method
		enforcer := model.GetCasbin()
		res1, err1 := enforcer.Enforce(sub1, obj, act)
		res2, err2 := enforcer.Enforce(sub2, obj, act)

		if err1 != nil || err2 != nil {
			utils.Consolelog(err1, err2)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "内部错误",
			})
			c.Abort()
			return
		} else if res1 || res2 {
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "很抱歉您没有此权限",
			})
			c.Abort()
		}
		c.Next()
	}
}
