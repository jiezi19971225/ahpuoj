package mysql

import (
	"log"
	"strings"
	"time"
	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Pool *sqlx.DB

func init() {
	var err error
	configFilePath := "config.ini"
	cfg, _ := goconfig.LoadConfigFile(configFilePath)
	dbcfg, _ := cfg.GetSection("mysql")
	path := strings.Join([]string{dbcfg["user"], ":", dbcfg["password"], "@tcp(", dbcfg["host"], ":", dbcfg["port"], ")/", dbcfg["database"], "?charset=utf8"}, "")
	Pool, err = sqlx.Open("mysql", path)
	Pool.SetMaxIdleConns(100)
	Pool.SetConnMaxLifetime(2 * time.Minute)
	if err != nil {
		log.Println("m=GetPool,msg=connection has failed", err)
	}
}
