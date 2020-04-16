package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/service/rabbitmq"
	"ahpuoj/utils"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*
前台用户登陆后可访问api的控制器
*/

// 用户获取用户信息
func GetUser(c *gin.Context) {

	user, _ := c.Get("user")
	if user, ok := user.(model.User); ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户信息获取成功",
			"user":    user,
		})
	}
}

// 账号设置中重设昵称的接口
func ResetNick(c *gin.Context) {
	var user model.User
	user, _ = GetUserInstance(c)
	var req struct {
		Nick string `json:"nick" binding:"required,max=20"`
	}
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}
	user.Nick = req.Nick
	_, err = DB.Exec("update user set nick = ? where id = ?", req.Nick, user.Id)
	if utils.CheckError(c, err, "该昵称已被使用") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "昵称修改成功",
		"show":    true,
		"user":    user,
	})
}

// 账号设置中重设密码的接口
func ResetPassword(c *gin.Context) {
	var user model.User
	user, _ = GetUserInstance(c)
	var req struct {
		OldPassword     string `json:"oldpassword" binding:"required,ascii,min=6,max=20"`
		Password        string `json:"password" binding:"required,ascii,min=6,max=20"`
		ConfirmPassword string `json:"confirmpassword" binding:"required,ascii,min=6,max=20"`
	}
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}
	h := sha1.New()
	h.Write([]byte(user.PassSalt))
	h.Write([]byte(req.OldPassword))
	hashedOldPassword := fmt.Sprintf("%x", h.Sum(nil))

	if hashedOldPassword != user.Password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "原密码错误",
			"show":    true,
		})
		return
	}

	// 更新密码
	// 加盐处理 16位随机字符串
	salt := utils.GetRandomString(16)
	h.Reset()
	h.Write([]byte(salt))
	h.Write([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	if hashedPassword == user.Password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "新密码不能和原密码相同",
			"show":    true,
		})
		return
	}

	_, err = DB.Exec("update user set password = ?, passsalt = ? where id = ?", hashedPassword, salt, user.Id)
	utils.Consolelog(err)
	c.JSON(http.StatusOK, gin.H{
		"message": "密码修改成功",
		"show":    true,
	})
}

// 用户提交测试运行的接口
func SubmitToTestRun(c *gin.Context) {
	var err error
	var req struct {
		Language  int    `json:"language" binding:"gte=0,lte=17"`
		InputText string `json:"input_text"  binding:"max=65535"`
		Source    string `json:"source"  binding:"required,min=2,max=65535"`
	}
	err = c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "提交失败,表单参数错误") != nil {
		return
	}

	conn := REDISPOOL.Get()
	defer conn.Close()

	testrunCount, _ := redis.Int(conn.Do("incr", "testrun:count"))
	testrunCountStr := strconv.Itoa(testrunCount)
	jsondata, _ := json.Marshal(gin.H{
		"UserId":       0,
		"TestrunCount": testrunCountStr,
		"SolutionId":   0,
		"ProblemId":    0,
		"Language":     req.Language,
		"TimeLimit":    1,
		"MemoryLimit":  64,
		"Source":       req.Source,
		"InputText":    req.InputText,
	})
	rabbitmq.Publish("oj", "problem", jsondata)
	//等待评测机评判
	var reinfo, ceinfo, costomOut string
	queryTimes := 0
	for {
		queryTimes += 1
		if queryTimes > 100 {
			break
		}
		exist, _ := redis.Int(conn.Do("exists", "testrun:"+testrunCountStr))
		if exist == 1 {
			values, _ := redis.Values(conn.Do("hmget", "testrun:"+testrunCountStr, "ceinfo", "reinfo", "custom_out"))
			redis.Scan(values, &reinfo, &ceinfo, &costomOut)
		} else {
			time.Sleep(time.Second)
		}
	}
	customOutput := reinfo
	if len(ceinfo) > 0 {
		customOutput = ceinfo
	}
	if len(costomOut) > 0 {
		customOutput = costomOut
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "测试运行成功",
		"show":          true,
		"custom_output": customOutput,
	})
}

// 用户提交评测的接口
func SubmitToJudge(c *gin.Context) {
	var err error
	var contest model.Contest
	user, _ := GetUserInstance(c)
	var req struct {
		ProblemId int    `json:"problem_id" binding:"required"`
		Language  int    `json:"language" binding:"gte=0,lte=17"`
		ContestId int    `json:"contest_id"`
		Num       int    `json:"num" binding:"omitempty,gte=0"`
		Source    string `json:"source"  binding:"required,min=2,max=65535"`
	}

	err = c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "提交失败") != nil {
		return
	}

	var problem model.Problem
	err = DB.Get(&problem, "select * from problem where id = ?", req.ProblemId)
	if utils.CheckError(c, err, "提交失败，问题不存在") != nil {
		return
	}

	submitable := false

	// 管理员无提交限制
	if user.Role != "user" {
		submitable = true
	} else {
		// 比赛的提交
		if req.ContestId > 0 {

			err = DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0", req.ContestId)
			if utils.CheckError(c, err, "提交失败，竞赛不存在") != nil {
				return
			}
			// 非管理员只有在比赛进行过程中并且有参加权限才能提交
			contest.CalcStatus()
			// 比赛进行中
			if contest.Status == 2 {
				// 公开
				if contest.Private == 0 {
					submitable = true
				} else {
					// 检测是否有提交权限
					var temp int
					DB.Get(&temp, "select count(1) from contest_user where contest_id = ? and user_id = ?", req.ContestId, user.Id)
					if temp > 0 {
						submitable = true
					}
				}
			}
		} else { // 平时的提交
			// 如果只是一般用户无法提交保留中的题目
			if problem.Defunct == 0 {
				submitable = true
			}
		}
	}

	if submitable {
		var teamId int
		// 如果为团队赛模式，并且非管理员提交，查询当前用户的teamId
		if contest.TeamMode == 1 && user.Role != "admin" {
			err = DB.Get(&teamId, "select team_id from contest_team_user ctu where ctu.contest_id = ? and ctu.user_id = ?", contest.Id, user.Id)
		}
		solution := model.Solution{
			ProblemId:  req.ProblemId,
			TeamId:     teamId,
			UserId:     user.Id,
			ContestId:  req.ContestId,
			Num:        req.Num,
			IP:         c.ClientIP(),
			Language:   req.Language,
			CodeLength: len(req.Source),
		}
		err := solution.Save()
		if utils.CheckError(c, err, "保存提交记录失败") != nil {
			return
		}
		sourceCode := model.SourceCode{
			SolutionId: solution.Id,
			Source:     req.Source,
		}
		err = sourceCode.Save()
		if utils.CheckError(c, err, "保存代码记录失败") != nil {
			return
		}
		// 将判题任务推入消息队列
		jsondata, _ := json.Marshal(gin.H{
			"UserId":       user.Id,
			"TestrunCount": 0,
			"SolutionId":   solution.Id,
			"ProblemId":    req.ProblemId,
			"Language":     req.Language,
			"TimeLimit":    problem.TimeLimit,
			"MemoryLimit":  problem.MemoryLimit,
			"Source":       req.Source,
			"InputText":    "",
		})

		rabbitmq.Publish("oj", "problem", jsondata)
		c.JSON(http.StatusOK, gin.H{
			"message":  "提交成功",
			"show":     true,
			"solution": solution,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "对不起，你没有提交权限",
			"show":    true,
		})
	}
}

// 切换代码公开状态
func ToggleSolutionStatus(c *gin.Context) {
	var user model.User
	id, _ := strconv.Atoi(c.Param("id"))

	user, _ = GetUserInstance(c)

	var solutionUserId int
	DB.Get(&solutionUserId, "select user_id from solution where solution_id = ?", id)
	if user.Id != solutionUserId {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "对不起，你没有修改权限",
		})
		return
	}

	DB.Exec("update source_code set public = not public where solution_id = ?", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "修改代码公开状态成功",
		"show":    true,
	})
}

// 下载题目数据文件
func DownloadDataFile(c *gin.Context) {
	var user model.User
	user, _ = GetUserInstance(c)

	pidStr := c.Query("pid")
	sidStr := c.Query("sid")
	filename := c.Query("filename")
	pid, _ := strconv.Atoi(pidStr)
	sid, _ := strconv.Atoi(sidStr)

	// 检验提交是否存在
	userId := 0
	err := DB.Get(&userId, "select 1 from solution where solution_id = ? and problem_id = ? and user_id = ?", sid, pid, user.Id)
	if utils.CheckError(c, err, "数据不存在") != nil {
		return
	}
	// 检验错误信息是否与数据库信息匹配
	errFilename := ""
	err = DB.Get(&errFilename, "select error from runtimeinfo where solution_id = ?", sid)
	if utils.CheckError(c, err, "数据不存在") != nil {
		return
	}
	errFilenameWithoutSuffix := strings.TrimSuffix(errFilename, filepath.Ext(errFilename))
	filenameWithoutSuffix := strings.TrimSuffix(filename, filepath.Ext(filename))
	if errFilenameWithoutSuffix != filenameWithoutSuffix {
		utils.Consolelog(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": "数据不存在",
				"show":    true,
			})
		return
	}

	// 读取文件
	cfg := utils.GetCfg()
	dataDir, _ := cfg.GetValue("project", "datadir")
	baseDir := dataDir + "/" + strconv.FormatInt(int64(pid), 10)
	dataFileName := baseDir + "/" + filename

	c.Header("Content-Disposition", `attachment; filename=`+filename+".txt")
	c.Header("Content-Type", "application/octet-stream")
	c.File(dataFileName)
}

// 账号设置上传头像
func UploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("img")
	ext := path.Ext(header.Filename)
	if utils.CheckError(c, err, "头像上传失败，参数错误") != nil {
		return
	}
	url, err := utils.SaveFile(file, ext, "avatars")

	if utils.CheckError(c, err, "头像上传失败,请检查服务器设置") != nil {
		return
	}

	var user model.User
	user, _ = GetUserInstance(c)
	// 如果不是默认头像 删除原头像
	defaultAvatar, _ := utils.GetCfg().GetValue("preset", "avatar")
	if user.Avatar != defaultAvatar {
		cfg := utils.GetCfg()
		webDir, _ := cfg.GetValue("project", "webdir")
		projectPath := webDir + "/"
		os.Remove(projectPath + user.Avatar)
	}
	DB.Exec("update user set avatar = ? where id = ?", url, user.Id)

	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"show":    true,
		"url":     url,
	})
}

// 发布主题帖
func PostIssue(c *gin.Context) {
	var err error
	var user model.User

	user, _ = GetUserInstance(c)
	var req request.Issue
	err = c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}
	if req.ProblemId != 0 {
		var temp int
		DB.Get(&temp, "select count(1) from problem where id = ?", req.ProblemId)
		if temp == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "发布讨论主题失败，问题不存在",
				"show":    true,
			})
			return
		}
	}
	issue := model.Issue{
		Title:     req.Title,
		ProblemId: req.ProblemId,
		UserId:    user.Id,
	}
	err = issue.Save()
	if utils.CheckError(c, err, "发布讨论主题失败") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "发布讨论主题成功",
		"show":    true,
		"issue":   issue,
	})
}

// 回复主题帖
func ReplyToIssue(c *gin.Context) {
	var err error
	var user model.User
	issueId, _ := strconv.Atoi(c.Param("id"))

	user, _ = GetUserInstance(c)
	var req request.Reply
	err = c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}
	// 主题是否存在
	var temp int
	DB.Get(&temp, "select count(1) from issue where id = ?", issueId)
	if temp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "发布回复失败，目标主题不存在",
			"show":    true,
		})
		return
	}
	// 如果是对回复的回复 检查该回复是否存在
	if req.ReplyId != 0 {
		DB.Get(&temp, "select count(1) from reply where id = ?", req.ReplyId)
		if temp == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "发布回复失败，目标回复不存在",
				"show":    true,
			})
			return
		}
	}
	reply := model.Reply{
		IssueId:     issueId,
		UserId:      user.Id,
		ReplyId:     req.ReplyId,
		ReplyUserId: req.ReplyUserId,
		Content:     req.Content,
	}
	err = reply.Save()
	if utils.CheckError(c, err, "发布回复失败，数据库操作错误") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "发布回复成功",
		"reply":   reply,
		"show":    true,
	})
}

// 获取回复我的信息帖子列表
func GetMyReplys(c *gin.Context) {
	var err error
	var user model.User
	user, _ = GetUserInstance(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perpage, _ := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	if page == 0 {
		page = 1
	}
	// 第一步只获取对主题的回复
	whereString := "where reply.reply_user_id = " + strconv.Itoa(user.Id)
	// + " and reply.user_id != " + strconv.Itoa(user.Id)
	// 管理员可以查看被删除的回复
	if user.Role != "admin" {
		whereString += " and reply.is_deleted = 0 "
	}
	rows, total, err := model.Paginate(&page, &perpage, "reply inner join user on reply.user_id = user.id inner join issue on reply.issue_id = issue.id",
		[]string{"user.username,user.nick,user.avatar,reply.*,'' as rnick,(select count(1) from reply  r where reply.id = r.reply_id) as reply_count,issue.title as issue_title"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	replys := []model.Reply{}
	for rows.Next() {
		var reply model.Reply
		rows.StructScan(&reply)
		replys = append(replys, reply)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"page":    page,
		"replys":  replys,
	})

}

// 获取最近一次提交的代码
func GetLatestSource(c *gin.Context) {
	var err error
	var user model.User
	user, _ = GetUserInstance(c)
	problemId, _ := strconv.Atoi(c.Param("id"))
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}
	// 查找提交
	type SourceCode struct {
		Source   string `json:"source"`
		Language int    `json:"language"`
	}
	var sourceCode SourceCode
	err = DB.Get(&sourceCode, `select source_code.source,solution.language from solution inner join source_code 
	on source_code.solution_id = solution.solution_id where solution.problem_id = ? and solution.user_id = ? order by solution.in_date desc limit 1`, problemId, user.Id)
	c.JSON(http.StatusOK, gin.H{
		"message":    "获取最近提交信息成功",
		"sourcecode": sourceCode,
	})
}

// 获取最近一次比赛中提交的代码
func GetLatestContestSource(c *gin.Context) {
	var err error
	var user model.User

	user, _ = GetUserInstance(c)
	contestId, _ := strconv.Atoi(c.Param("id"))
	num, _ := strconv.Atoi(c.Param("num"))

	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}

	// 查找提交
	type SourceCode struct {
		Source   string `json:"source"`
		Language int    `json:"language"`
	}
	var sourceCode SourceCode
	err = DB.Get(&sourceCode, `select source_code.source,solution.language from solution inner join source_code 
	on source_code.solution_id = solution.solution_id where solution.contest_id = ? and solution.num = ? and solution.user_id = ? order by solution.in_date desc limit 1`, contestId, num, user.Id)
	c.JSON(http.StatusOK, gin.H{
		"message":    "获取最近提交信息成功",
		"sourcecode": sourceCode,
	})
}
