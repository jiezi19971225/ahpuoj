package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTeam(c *gin.Context) {
	var team model.Team
	id, _ := strconv.Atoi(c.Param("id"))
	err := DB.Get(&team, "select * from team where id = ?", id)
	if utils.CheckError(c, err, "团队不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"team":    team,
	})
}

func IndexTeam(c *gin.Context) {
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := " where is_deleted = 0 "
	if len(param) > 0 {
		whereString += "and name like '%" + param + "%'"
	}
	whereString += " order by team.id desc"
	rows, total, err := model.Paginate(&page, &perpage, "team inner join user on team.user_id = user.id", []string{"team.*", "user.username"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	teams := []model.Team{}
	for rows.Next() {
		var team model.Team
		rows.StructScan(&team)
		teams = append(teams, team)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"page":    page,
		"perpage": perpage,
		"data":    teams,
	})
}

func GetAllTeams(c *gin.Context) {
	rows, _ := DB.Queryx("select * from team where is_deleted = 0 order by id desc")
	teams := []model.Team{}
	for rows.Next() {
		var team model.Team
		rows.StructScan(&team)
		teams = append(teams, team)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"teams":   teams,
	})
}

func IndexTeamUser(c *gin.Context) {
	teamId := c.Param("id")
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := "where team_user.team_id=" + teamId
	if len(param) > 0 {
		whereString += " and user.username like '%" + param + "%' or user.nick like '%" + param + "%'"
	}
	whereString += " order by user.id desc"
	rows, total, err := model.Paginate(&page, &perpage,
		"team_user inner join user on team_user.user_id = user.id",
		[]string{"user.*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	users := []model.User{}
	for rows.Next() {
		var user model.User
		rows.StructScan(&user)
		users = append(users, user)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"page":    page,
		"perpage": perpage,
		"data":    users,
	})
}

func AddTeamUsers(c *gin.Context) {
	var req struct {
		UserList string `json:"userlist" binding:"required"`
	}
	id, _ := strconv.Atoi(c.Param("id"))
	c.ShouldBindJSON(&req)
	team := model.Team{
		Id: id,
	}
	infos := team.AddUsers(req.UserList)
	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"show":    true,
		"info":    infos,
	})
}

func StoreTeam(c *gin.Context) {
	var req request.Team
	user, _ := GetUserInstance(c)
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	team := model.Team{
		Name:   req.Name,
		UserId: user.Id,
	}
	err = team.Save()
	// TODO 先这样处理 给次级管理员添加权限
	if utils.CheckError(c, err, "新建团队失败，该团队已存在") != nil {
		return
	}
	idStr := strconv.Itoa(user.Id)
	teamIdStr := strconv.Itoa(team.Id)
	if user.Role != "admin" {
		enforcer := model.GetCasbin()
		enforcer.AddPolicy(idStr, "/api/admin/team/"+teamIdStr+"/users", "POST")
		enforcer.AddPolicy(idStr, "/api/admin/team/"+teamIdStr, "PUT")
		enforcer.AddPolicy(idStr, "/api/admin/team/"+teamIdStr+"/user/:userid", "DELETE")
		enforcer.AddPolicy(idStr, "/api/admin/team/"+teamIdStr, "DELETE")
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "新建团队成功",
		"show":    true,
		"team":    team,
	})
}

func UpdateTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.Team
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	team := model.Team{
		Id:   id,
		Name: req.Name,
	}
	err = team.Update()
	if utils.CheckError(c, err, "编辑团队失败") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "编辑团队成功",
		"show":    true,
		"team":    team,
	})
}

func DeleteTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	team := model.Team{
		Id: id,
	}
	err := team.Delete()
	if utils.CheckError(c, err, "删除团队失败，团队不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除团队成功",
		"show":    true,
	})
}

func DeleteTeamUser(c *gin.Context) {
	teamId, _ := strconv.Atoi(c.Param("id"))
	userId, _ := strconv.Atoi(c.Param("userid"))
	result, _ := DB.Exec("delete from team_user where team_id = ? and user_id = ?", teamId, userId)
	// 级联删除
	DB.Exec(`delete contest_user from contest_user inner join contest_team_user on contest_user.contest_id = contest_team_user.contest_id 
	where contest_user.user_id = ? and contest_team_user.team_id = ?`, userId, teamId)
	DB.Exec("delete from contest_team_user where team_id = ? and user_id = ?", teamId, userId)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除团队成员失败，团队成员不存在",
			"show":    true,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除团队成员成功",
		"show":    true,
	})
}
