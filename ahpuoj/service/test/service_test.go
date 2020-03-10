package test

import (
	"ahpuoj/utils"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/gomodule/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"github.com/streadway/amqp"
)

func TestDB(t *testing.T) {

	cfg := utils.GetTestCfg("../../config/config.ini")
	dbcfg, _ := cfg.GetSection("mysql")
	path := strings.Join([]string{dbcfg["user"], ":", dbcfg["password"], "@tcp(", dbcfg["host"], ":", dbcfg["port"], ")/", dbcfg["database"], "?charset=utf8"}, "")
	t.Log(path)
	POOL, err := sqlx.Open("mysql", path)
	err = POOL.Ping()
	if err != nil {
		t.Error("m=GetPool,msg=connection has failed", err)
	}
}

func TestRedis(t *testing.T) {

	cfg := utils.GetTestCfg("../../config/config.ini")
	rediscfg, _ := cfg.GetSection("redis")
	path := strings.Join([]string{rediscfg["host"], ":", rediscfg["port"]}, "")
	t.Log(path)
	Pool := &redis.Pool{
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
	conn := Pool.Get()
	err := Pool.TestOnBorrow(conn, time.Now())
	if err != nil {
		t.Error("connect redis failed", err)
	}
}

func TestMQ(t *testing.T) {
	cfg := utils.GetTestCfg("../../config/config.ini")
	mqcfg, _ := cfg.GetSection("rabbitmq")
	path := strings.Join([]string{"amqp://", mqcfg["user"], ":", mqcfg["password"], "@", mqcfg["host"], ":", mqcfg["port"], "/oj"}, "")
	t.Log(path)
	_, err := amqp.Dial(path)
	if err != nil {
		t.Error("connect mq failed", err)
	}
}
