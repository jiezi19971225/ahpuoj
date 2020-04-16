package model

import (
	"ahpuoj/utils"
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"strings"
	//xormadapter "github.com/casbin/xorm-adapter"
)

type Casbin struct {
	ID       int    `json:"id" `
	Ptype    string `json:"ptype"`
	RoleName string `json:"rolename"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}

func GetCasbin() *casbin.Enforcer {
	cfg := utils.GetCfg()
	dbcfg, _ := cfg.GetSection("mysql")
	path := strings.Join([]string{dbcfg["user"], ":", dbcfg["password"], "@tcp(", dbcfg["host"], ":", dbcfg["port"], ")/"}, "")
	adapter, _ := xormadapter.NewAdapter("mysql", path)
	configFilePath := "config/auth_model.conf"
	enforcer, _ := casbin.NewEnforcer(configFilePath, adapter)
	enforcer.LoadPolicy()
	return enforcer
}

func (c *Casbin) Store() error {
	enforcer := GetCasbin()
	enforcer.AddPolicy(c.RoleName, c.Path, c.Method)
	return enforcer.SavePolicy()
}
