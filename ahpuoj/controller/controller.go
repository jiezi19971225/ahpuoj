package controller

import (
	"ahpuoj/model"
	"ahpuoj/service/mysql"
	"ahpuoj/service/redisConn"
	"ahpuoj/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB
var REDISPOOL *redis.Pool
var RedisCacheLiveTime int

func init() {
	DB = mysql.Pool
	REDISPOOL = redisConn.Pool

	// 默认1800
	RedisCacheLiveTime = 1800
	if rcltstr, err := utils.GetCfg().GetValue("redis", "cacheLiveTime"); err == nil {
		if rclt, err := strconv.Atoi(rcltstr); err == nil {
			RedisCacheLiveTime = rclt
		}
	}
}

// 获得user实例
func GetUserInstance(c *gin.Context) (model.User, bool) {
	var user model.User
	userInterface, loggedIn := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}
	return user, loggedIn
}
