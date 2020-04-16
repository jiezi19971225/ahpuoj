package controller

import (
	"ahpuoj/model"
	"ahpuoj/utils"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexUser(c *gin.Context) {

	userType := c.Query("userType")
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := " where is_compete_user =" + userType
	if len(param) > 0 {
		whereString += " and username like '%" + param + "%' or nick like '%" + param + "%'"
	}
	whereString += " order by id desc"
	rows, total, err := model.Paginate(&page, &perpage, "user", []string{"*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	users := []model.User{}
	for rows.Next() {
		var user model.User
		rows.StructScan(&user)
		users = append(users, user)
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"page":    page,
		"perpage": perpage,
		"data":    users,
	})
}

func ToggleUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := model.User{
		Id: id,
	}
	err := user.ToggleStatus()
	if utils.CheckError(c, err, "更改用户状态失败，用户不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改用户状态成功",
		"show":    true,
	})
}

func ChangeUserPass(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Password string `json:"password" binding:"required,ascii,min=6,max=20"`
	}
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	// 更新密码
	// 加盐处理 16位随机字符串
	salt := utils.GetRandomString(16)
	h := sha1.New()
	h.Write([]byte(salt))
	h.Write([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	user := model.User{
		Id:       id,
		Password: hashedPassword,
		PassSalt: salt,
	}
	err = user.ChangePass()
	if utils.CheckError(c, err, "更改用户密码失败，用户不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改用户密码成功",
		"show":    true,
	})
}
