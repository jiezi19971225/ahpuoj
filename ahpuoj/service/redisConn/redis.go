package redisConn

import (
	"ahpuoj/utils"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

var Pool *redis.Pool

func init() {
	cfg := utils.GetCfg()
	rediscfg, _ := cfg.GetSection("redis")
	path := strings.Join([]string{rediscfg["host"], ":", rediscfg["port"]}, "")

	Pool = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   10000,
		Wait:        true,
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", path)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", rediscfg["password"]); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
