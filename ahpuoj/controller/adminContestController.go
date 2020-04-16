package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"archive/zip"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func IndexContest(c *gin.Context) {
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := " where is_deleted = 0 "
	if len(param) > 0 {
		whereString += "and name like '%" + param + "%'"
	}
	whereString += " order by contest.id desc"
	rows, total, err := model.Paginate(&page, &perpage, "contest inner join user on contest.user_id = user.id", []string{"contest.*,user.username"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	contests := []model.Contest{}
	for rows.Next() {
		var contest model.Contest
		rows.StructScan(&contest)
		contests = append(contests, contest)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"page":    page,
		"perpage": perpage,
		"data":    contests,
	})
}

func ShowContest(c *gin.Context) {
	var contest model.Contest
	id, _ := strconv.Atoi(c.Param("id"))
	err := DB.Get(&contest, "select * from contest where id = ?", id)
	if utils.CheckError(c, err, "竞赛&作业不存在") != nil {
		return
	}
	contest.FetchProblems()
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"contest": contest.Response(),
	})
}

func GetAllContests(c *gin.Context) {
	rows, _ := DB.Queryx("select * from contest where is_deleted = 0")
	var contests []map[string]interface{}
	for rows.Next() {
		var contest model.Contest
		rows.StructScan(&contest)
		contests = append(contests, contest.ListItemResponse())
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "数据获取成功",
		"contests": contests,
	})
}

func StoreContest(c *gin.Context) {
	var req request.Contest
	user, _ := GetUserInstance(c)
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	contest := model.Contest{
		Name:        req.Name,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Description: model.NullString{sql.NullString{req.Description, true}},
		LangMask:    req.LangMask,
		Private:     req.Private,
		TeamMode:    req.TeamMode,
		UserId:      user.Id,
	}
	err = contest.Save()
	// 处理竞赛作业包含的问题
	contest.AddProblems(req.Problems)
	if utils.CheckError(c, err, "新建竞赛&作业失败，该竞赛&作业已存在") != nil {
		return
	}
	idStr := strconv.Itoa(user.Id)
	contestIdStr := strconv.Itoa(contest.Id)
	if user.Role != "admin" {
		enforcer := model.GetCasbin()
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr, "PUT")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr, "DELETE")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/status", "PUT")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/users", "POST")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/user/:userid", "DELETE")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/team/:teamid", "POST")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/team/:teamid", "DELETE")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/team/:teamid/users", "POST")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/team/:teamid/allusers", "POST")
		enforcer.AddPolicy(idStr, "/api/admin/contest/"+contestIdStr+"/team/:teamid/user/:userid", "DELETE")
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "新建竞赛&作业成功",
		"show":    true,
		"contest": contest,
	})
}

func UpdateContest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.Contest
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	contest := model.Contest{
		Id:          id,
		Name:        req.Name,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Description: model.NullString{sql.NullString{req.Description, true}},
		LangMask:    req.LangMask,
		Private:     req.Private,
		TeamMode:    req.TeamMode,
	}
	err = contest.Update()
	if utils.CheckError(c, err, "编辑竞赛&作业失败，竞赛&作业不存在") != nil {
		return
	}
	// 处理题目列表
	contest.RemoveProblems()
	contest.AddProblems(req.Problems)
	c.JSON(http.StatusOK, gin.H{
		"message": "编辑竞赛&作业成功",
		"show":    true,
		"contest": contest,
	})
}

func DeleteContest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contest := model.Contest{
		Id: id,
	}
	err := contest.Delete()
	if utils.CheckError(c, err, "删除竞赛&作业失败，竞赛&作业不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除竞赛&作业成功",
		"show":    true,
	})
}

func ToggleContestStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contest := model.Contest{
		Id: id,
	}
	err := contest.ToggleStatus()
	if utils.CheckError(c, err, "更改竞赛&作业状态失败，竞赛&作业不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改竞赛&作业状态成功",
		"show":    true,
	})
}

// 处理个人赛人员列表
func IndexContestUser(c *gin.Context) {
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := "where contest_user.contest_id=" + c.Param("id")
	if len(param) > 0 {
		whereString += " and user.username like '%" + param + "%' or user.nick like '%" + param + "%'"
	}
	whereString += " order by user.id desc"
	rows, total, err := model.Paginate(&page, &perpage,
		"contest_user inner join user on contest_user.user_id = user.id",
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

func AddContestUsers(c *gin.Context) {
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		UserList string `json:"userlist" binding:"required"`
	}
	c.ShouldBindJSON(&req)
	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}
	contest := model.Contest{
		Id: id,
	}
	infos := contest.AddUsers(req.UserList, 0)
	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"show":    true,
		"info":    infos,
	})
}

func DeleteContestUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userId, _ := strconv.Atoi(c.Param("userid"))
	DB.Exec("delete from contest_user where contest_id = ? and user_id = ?", id, userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除竞赛&作业人员成功",
		"show":    true,
	})
}

// 处理团队赛管理
func IndexContestTeamWithUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rows, _ := DB.Queryx("select team.* from contest_team inner join team on contest_team.team_id = team.id where contest_team.contest_id = ?", id)
	teams := []model.Team{}
	for rows.Next() {
		var team model.Team
		rows.StructScan(&team)
		team.AttachUserInfo(id)
		teams = append(teams, team)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"data":    teams,
	})
}

func AddContestTeam(c *gin.Context) {
	var err error
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}
	// 检查团队是否存在
	DB.Get(&temp, "select count(1) from team where id = ? and is_deleted = 0", teamId)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "团队不存在",
		})
		return
	}
	// 检查是否已经添加进了竞赛作业中
	DB.Get(&temp, "select count(1) from contest_team where contest_id = ? and team_id = ? ", id, teamId)
	if temp > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "该团队已经在该竞赛作业中",
		})
		return
	}
	_, err = DB.Exec("insert into contest_team(contest_id,team_id,created_at,updated_at) values(?,?,NOW(),NOW())", id, teamId)
	if utils.CheckError(c, err, "数据库操作失败") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "添加团队成功",
		"show":    true,
	})
}

func DeleteContestTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	DB.Exec("delete from contest_team where contest_id = ? and team_id = ?", id, teamId)
	// 级联删除
	DB.Exec(`delete contest_user from contest_user inner join contest_team_user on contest_user.contest_id = contest_team_user.contest_id 
	where contest_user.contest_id = ? and contest_team_user.team_id = ?`, id, teamId)
	DB.Exec("delete from contest_team_user where contest_id = ? and team_id = ?", id, teamId)

	c.JSON(http.StatusOK, gin.H{
		"message": "删除团队成功",
		"show":    true,
	})
}

func DeleteContestTeamUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	userId, _ := strconv.Atoi(c.Param("userid"))
	DB.Exec(`delete contest_user from contest_user inner join contest_team_user on contest_user.contest_id = contest_team_user.contest_id 
	where contest_user.contest_id = ? and contest_user.user_id = ? and contest_team_user.team_id = ?`, id, userId, teamId)
	// 级联删除
	DB.Exec("delete from contest_team_user where contest_id = ? and team_id = ? and user_id = ?", id, teamId, userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除团队人员成功",
		"show":    true,
	})
}

func AddContestTeamUsers(c *gin.Context) {
	var req struct {
		UserList string `json:"userlist" binding:"required"`
	}
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	c.ShouldBindJSON(&req)

	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}

	// 检查团队是否存在
	DB.Get(&temp, "select count(1) from team where id = ? and is_deleted = 0", teamId)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "团队不存在",
		})
		return
	}

	contest := model.Contest{
		Id: id,
	}
	infos := contest.AddUsers(req.UserList, teamId)
	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"show":    true,
		"info":    infos,
	})
}

func AddContestTeamAllUsers(c *gin.Context) {
	var err error
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}
	// 检查团队是否存在
	DB.Get(&temp, "select count(1) from team where id = ? and is_deleted = 0", teamId)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "团队不存在",
		})
		return
	}
	var infos []string
	rows, err := DB.Queryx("select user.* from user inner join team_user on user.id = team_user.user_id where team_user.team_id = ?", teamId)
	checkHasUserStmt, _ := DB.Preparex("select 1 from contest_user where contest_user.contest_id = ? and contest_user.user_id = ?")
	insertStmt, _ := DB.Preparex("insert into contest_user(contest_id,user_id,created_at,updated_at) VALUES (?,?,NOW(),NOW())")
	insertToTeamStmt, _ := DB.Preparex("insert into contest_team_user(contest_id,team_id,user_id,created_at,updated_at) VALUES (?,?,?,NOW(),NOW())")
	for rows.Next() {
		var user model.User
		var info string
		rows.StructScan(&user)
		utils.Consolelog(user)
		err = checkHasUserStmt.Get(&temp, id, user.Id)
		// 有记录返回err==nil
		if err == nil {
			info = "竞赛&作业添加用户" + user.Username + "失败，用户已被添加"
		} else {
			insertStmt.Exec(id, user.Id)
			insertToTeamStmt.Exec(id, teamId, user.Id)
			info = "竞赛&作业添加用户" + user.Username + "成功"
		}
		infos = append(infos, info)
	}
	insertStmt.Close()
	insertToTeamStmt.Close()
	checkHasUserStmt.Close()
	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"show":    true,
		"info":    infos,
	})

}

func GetContestProblemSolutions(c *gin.Context) {
	//var err error
	//var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	problemId, _ := strconv.Atoi(c.Param("problemid"))
	type SolutionWithName struct {
		Source     string `db:"source" json:"source"`
		Username   string `db:"username" json:"username"`
		SolutionId int    `db:"solution_id" json:"solution_id"`
		Language   int    `db:"language" json:"language"`
	}
	var SolutionWithNameList []SolutionWithName
	DB.Select(&SolutionWithNameList, "select source_code.solution_id,source_code.source,username,language from solution "+
		"INNER JOIN source_code on solution.solution_id = source_code.solution_id "+
		"INNER JOIN user on solution.user_id = user.id "+
		"where contest_id = ? and num = ? and result = 4 "+
		" order by username desc,solution_id desc;", id, problemId)

	buf := new(bytes.Buffer)
	// 实例化新的 zip.Writer
	w := zip.NewWriter(buf)

	// 从数据库查到的数据按照 用户名 提交id 排序，同一用户名可能有多份提交，需要手动去重
	prename := ""
	for _, solution := range SolutionWithNameList {
		if prename == solution.Username {
			continue
		}
		prename = solution.Username
		f, err := w.Create(solution.Username + "." + utils.LanguageExt[solution.Language])
		if err != nil {
			utils.Consolelog(err)
		}
		_, err = f.Write([]byte(solution.Source))
		if err != nil {
			utils.Consolelog(err)
		}
	}
	w.Close()

	filename := "c" + c.Param("id") + "t" + c.Param("problemid")
	contentLength := int64(buf.Len())

	// 异常处理
	defer func() {
		recover()
	}()
	c.DataFromReader(200, contentLength, `application/octet-stream`, buf, map[string]string{
		"Content-Disposition": `attachment; filename=` + filename + "solutions.zip",
	})
}
