package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRoleList(c *gin.Context) {
	type role struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	rolelist := []role{}
	DB.Unsafe().Select(&rolelist, "select * from role where name != 'user' and name != 'admin'")
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"data":    rolelist,
	})
}

func GetAdminList(c *gin.Context) {
	adminList := []model.User{}
	DB.Unsafe().Select(&adminList, "select user.*,role.name as role from user inner join role on user.role_id = role.id where user.role_id != 1")
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"data":    adminList,
	})
}

func AssignRole(c *gin.Context) {
	var req request.AssignRole
	err := c.ShouldBindJSON(&req)
	// TODO 先粗略处理 无法授予超级管理员角色
	msg := "请求参数错误"
	if req.RoleId == 2 {
		err = errors.New("无法授予超级管理员角色")
		msg = "无法授予超级管理员角色"
	}
	if utils.CheckError(c, err, msg) != nil {
		return
	}
	var temp int
	err = DB.Get(&temp, "select 1 from user where  username = ?", req.UserName)
	if utils.CheckError(c, err, "用户不存在") != nil {
		return
	}
	err = DB.Get(&temp, "select 1 from role where id = ?", req.RoleId)
	if utils.CheckError(c, err, "角色不存在") != nil {
		return
	}
	DB.Exec("update user set role_id = ? where username = ?", req.RoleId, req.UserName)
	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"show":    true,
	})
}

func UnassignRole(c *gin.Context) {
	var req request.UnassignRole
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	var user model.User
	err = DB.Unsafe().Get(&user, "select * from user where id = ?", req.UserId)
	// TODO 先粗略处理 无法撤销超级管理员角色
	msg := "用户不存在"
	if user.RoleId == 2 {
		err = errors.New("无法撤销超级管理员角色")
		msg = "无法撤销超级管理员角色"
	}
	if utils.CheckError(c, err, msg) != nil {
		return
	}
	DB.Exec("update user set role_id = 1 where id = ?", req.UserId)
	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"show":    true,
	})
}
