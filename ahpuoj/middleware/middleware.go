package middleware

import (
	"ahpuoj/service/redisConn"

	"github.com/gomodule/redigo/redis"
)

var REDISPOOL *redis.Pool

func init() {
	REDISPOOL = redisConn.Pool
}
