package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexSeries(c *gin.Context) {
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := " where is_deleted = 0 "
	if len(param) > 0 {
		whereString += "and name like '%" + param + "%'"
	}
	whereString += " order by series.id desc"
	rows, total, err := model.Paginate(&page, &perpage, "series inner join user on series.user_id = user.id", []string{"series.*,user.username"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	serieses := []model.Series{}
	for rows.Next() {
		var series model.Series
		err = rows.StructScan(&series)
		serieses = append(serieses, series)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"page":    page,
		"perpage": perpage,
		"data":    serieses,
	})
}

func ShowSeries(c *gin.Context) {
	var series model.Series
	id, _ := strconv.Atoi(c.Param("id"))
	err := DB.Get(&series, "select * from series where id = ?", id)
	if utils.CheckError(c, err, "系列赛不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"series":  series,
	})
}

func IndexSeriesContest(c *gin.Context) {
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := "where contest_series.series_id=" + c.Param("id")
	if len(param) > 0 {
		whereString += " and contest.name like '%" + param + "%'"
	}
	whereString += " order by contest.id desc"
	rows, total, err := model.Paginate(&page, &perpage,
		"contest_series inner join contest on contest_series.contest_id = contest.id",
		[]string{"contest.*"}, whereString)
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

func StoreSeries(c *gin.Context) {
	var req request.Series
	err := c.ShouldBindJSON(&req)
	user, _ := GetUserInstance(c)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	series := model.Series{
		Name:        req.Name,
		Description: req.Description,
		TeamMode:    req.TeamMode,
		UserId:      user.Id,
	}
	err = series.Save()
	if utils.CheckError(c, err, "新建系列赛失败，该系列赛已存在") != nil {
		return
	}
	idStr := strconv.Itoa(user.Id)
	seriesIdStr := strconv.Itoa(series.Id)
	if user.Role != "admin" {
		enforcer := model.GetCasbin()
		enforcer.AddPolicy(idStr, "/api/admin/series/"+seriesIdStr, "PUT")
		enforcer.AddPolicy(idStr, "/api/admin/series/"+seriesIdStr, "DELETE")
		enforcer.AddPolicy(idStr, "/api/admin/series/"+seriesIdStr+"/status", "PUT")
		enforcer.AddPolicy(idStr, "/api/admin/series/"+seriesIdStr+"/contest/:contestid", "POST")
		enforcer.AddPolicy(idStr, "/api/admin/series/"+seriesIdStr+"/contest/:contestid", "DELETE")
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "新建系列赛成功",
		"show":    true,
		"series":  series,
	})
}

func UpdateSeries(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.Series
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	series := model.Series{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		TeamMode:    req.TeamMode,
	}
	err = series.Update()
	if utils.CheckError(c, err, "编辑系列赛失败，该系列赛已存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "编辑系列赛成功",
		"show":    true,
		"series":  series,
	})
}

func ToggleSeriesStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	series := model.Series{
		Id: id,
	}
	err := series.ToggleStatus()
	if utils.CheckError(c, err, "更改系列赛状态失败，系列赛不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改系列赛状态成功",
		"show":    true,
	})
}

func DeleteSeries(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	series := model.Series{
		Id: id,
	}
	err := series.Delete()
	if utils.CheckError(c, err, "删除系列赛失败，该系列赛不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除系列赛成功",
		"show":    true,
	})
}

func AddSeriesContest(c *gin.Context) {
	var err error
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	contestId, _ := strconv.Atoi(c.Param("contestid"))
	// 检查系列赛是否存在
	DB.Get(&temp, "select count(1) from series where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "系列赛不存在",
			"show":    true,
		})
		return
	}

	// 检查竞赛&作业是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", contestId)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "竞赛&作业不存在",
			"show":    true,
		})
		return
	}

	// 检查是否已经添加进了竞赛作业中
	DB.Get(&temp, "select count(1) from contest_series where series_id = ? and contest_id = ? ", id, contestId)
	if temp > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "该竞赛&作业已经在系列赛中了",
			"show":    true,
		})
		return
	}
	_, err = DB.Exec("insert into contest_series(series_id,contest_id,created_at,updated_at) values(?,?,NOW(),NOW())", id, contestId)
	if utils.CheckError(c, err, "数据库操作失败") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "添加竞赛&作业成功",
		"show":    true,
	})
}

func DeleteSeriesContest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contestId, _ := strconv.Atoi(c.Param("contestid"))
	DB.Exec("delete from contest_series where series_id = ? and contest_id = ?", id, contestId)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除系列赛竞赛&作业成功",
		"show":    true,
	})
}
