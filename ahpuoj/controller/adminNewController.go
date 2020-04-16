package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexNew(c *gin.Context) {
	param := c.Query("param")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	whereString := ""
	if len(param) > 0 {
		whereString += "where title like '%" + param + "%'"
	}
	whereString += " order by top desc, id desc"
	rows, total, err := model.Paginate(&page, &perpage, "new", []string{"*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	news := []model.New{}
	for rows.Next() {
		var new model.New
		rows.StructScan(&new)
		news = append(news, new)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"page":    page,
		"perpage": perpage,
		"data":    news,
	})
}

func ShowNew(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	new := model.New{
		Id: id,
	}
	err := DB.Get(&new, "select * from new where id = ?", new.Id)
	if utils.CheckError(c, err, "新闻不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"new":     new,
	})
}

func StoreNew(c *gin.Context) {
	var req request.New
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	new := model.New{
		Title:   req.Title,
		Content: model.NullString{sql.NullString{req.Content, true}},
	}
	err = new.Save()
	if utils.CheckError(c, err, "新建新闻失败，该新闻已存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "新建新闻成功",
		"new":     new,
	})
}

func UpdateNew(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.New
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	new := model.New{
		Id:      id,
		Title:   req.Title,
		Content: model.NullString{sql.NullString{req.Content, true}},
	}
	err = new.Update()
	if utils.CheckError(c, err, "编辑新闻失败，该新闻已存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "编辑新闻成功",
		"new":     new,
	})
}

func DeleteNew(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	new := model.New{
		Id: id,
	}
	err := new.Delete()
	if utils.CheckError(c, err, "删除新闻失败，该新闻不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除新闻成功",
	})
}

func ToggleNewStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	new := model.New{
		Id: id,
	}
	err := new.ToggleStatus()
	if utils.CheckError(c, err, "更改新闻状态失败，该新闻不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改新闻状态成功",
	})
}

func ToggleNewTopStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	new := model.New{
		Id: id,
	}
	err := new.ToggleTopStatus()
	if utils.CheckError(c, err, "更改新闻置顶状态失败，该新闻不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改新闻置顶状态成功",
	})
}
