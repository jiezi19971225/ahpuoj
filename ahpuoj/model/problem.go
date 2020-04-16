package model

import (
	"ahpuoj/utils"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Problem struct {
	Id           int                      `db:"id" json:"id"`
	Title        string                   `db:"title" json:"title" binding:"required,max=20"`
	Description  NullString               `db:"description" json:"description"`
	Level        int                      `db:"level" json:"level" binding:"gte=0, lte=2"`
	Input        NullString               `db:"input" json:"input"`
	Output       NullString               `db:"output" json:"output"`
	SampleInput  NullString               `db:"sample_input" json:"sample_input"`
	SampleOutput NullString               `db:"sample_output" json:"sample_output"`
	Spj          int                      `db:"spj" json:"spj"`
	Hint         NullString               `db:"hint" json:"hint"`
	Defunct      int                      `db:"defunct" json:"defunct"`
	TimeLimit    int                      `db:"time_limit" json:"time_limit" binding:"required"`
	MemoryLimit  int                      `db:"memory_limit" json:"memory_limit" binding:"required"`
	Accepted     int                      `db:"accepted" json:"accepted"`
	Submit       int                      `db:"submit" json:"submit"`
	Solved       int                      `db:"solved" json:"solved"`
	CreatedAt    string                   `db:"created_at" json:"created_at"`
	UpdatedAt    string                   `db:"updated_at" json:"updated_at"`
	CreatorId    string                   `db:"creator_id" json:"creator_id"`
	Tags         []map[string]interface{} `json:"tags"`
	UserId       int                      `db:"user_id" json:"user_id"`
	Username     string                   `db:"username" json:"username"`
}

type ProblemWithoutTag struct {
	Id           int        `db:"id" json:"id"`
	Title        string     `db:"title" json:"title"`
	Description  NullString `db:"description" json:"description"`
	Level        int        `db:"level" json:"level"`
	Input        NullString `db:"input" json:"input"`
	Output       NullString `db:"output" json:"output"`
	SampleInput  NullString `db:"sample_input" json:"sample_input"`
	SampleOutput NullString `db:"sample_output" json:"sample_output"`
	Spj          int        `db:"spj" json:"spj"`
	Hint         NullString `db:"hint" json:"hint"`
	Defunct      int        `db:"defunct" json:"defunct"`
	TimeLimit    int        `db:"time_limit" json:"time_limit"`
	MemoryLimit  int        `db:"memory_limit" json:"memory_limit"`
	Accepted     int        `db:"accepted" json:"accepted"`
	Submit       int        `db:"submit" json:"submit"`
	Solved       int        `db:"solved" json:"solved"`
	CreatedAt    string     `db:"created_at" json:"created_at"`
	UpdatedAt    string     `db:"updated_at" json:"updated_at"`
	CreatorId    string     `db:"creator_id" json:"creator_id"`
	Tags         []Tag      `json:"-"`
}

func (problem *Problem) MarshalJSON() ([]byte, error) {
	type Alias Problem
	problem.Description.String = utils.ConvertTextImgUrl(problem.Description.String)
	problem.Input.String = utils.ConvertTextImgUrl(problem.Input.String)
	problem.Output.String = utils.ConvertTextImgUrl(problem.Output.String)
	problem.Hint.String = utils.ConvertTextImgUrl(problem.Hint.String)
	return json.Marshal((*Alias)(problem))
}

func (problem *Problem) Save() error {
	result, err := DB.Exec(`insert into problem
	(title,description,input,output,sample_input,sample_output,spj,hint,level,time_limit,memory_limit,defunct,user_id,created_at,updated_at) 
	values (?,?,?,?,?,?,?,?,?,?,?,?,?,NOW(),NOW())`, problem.Title, problem.Description, problem.Input, problem.Output,
		problem.SampleInput, problem.SampleOutput, 0, problem.Hint, problem.Level, problem.TimeLimit, problem.MemoryLimit, 1, problem.UserId)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	problem.Id = utils.Int64to32(lastInsertId)
	// 创建数据文件夹
	cfg := utils.GetCfg()
	dataDir, _ := cfg.GetValue("project", "datadir")
	baseDir := dataDir + "/" + strconv.FormatInt(int64(problem.Id), 10)
	utils.Consolelog(baseDir)
	err = os.MkdirAll(baseDir, 0777)
	return err
}

func (problem *Problem) Update() error {
	result, err := DB.Exec(`update problem 
	set title = ?,description = ?,input = ?,output = ?,sample_input = ?,sample_output = ?,spj = ?,
	hint = ?,level=?,time_limit = ?,memory_limit = ?,updated_at = NOW() where id = ?`, problem.Title, problem.Description, problem.Input, problem.Output,
		problem.SampleInput, problem.SampleOutput, 0, problem.Hint, problem.Level, problem.TimeLimit, problem.MemoryLimit, problem.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (problem *Problem) Delete() error {
	result, err := DB.Exec(`delete from problem where id = ?`, problem.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (problem *Problem) ToggleStatus() error {
	result, err := DB.Exec(`update problem set defunct = not defunct,updated_at = NOW() where id = ?`, problem.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (problem *Problem) AddTags(reqTags []interface{}) {

	tags := []map[string]interface{}{}
	insertStmt, _ := DB.Preparex(`insert into problem_tag (problem_id,tag_id,created_at,updated_at) values (?,?,NOW(),NOW())`)
	for _, tag := range reqTags {

		switch t := tag.(type) {
		// 如果是标签id 直接插入 golang解析json会将数字转成float64
		case float64:
			insertStmt.Exec(problem.Id, int(t))
			tag := Tag{
				Id: int(t),
			}
			tags = append(tags, gin.H{
				"id":   tag.Id,
				"name": tag.Name,
			})
			// 否则生成新的tag并且插入
		case string:
			newTag := Tag{
				Name: t,
			}
			err := newTag.Save()
			if err == nil {
				insertStmt.Exec(problem.Id, newTag.Id)
				tags = append(tags, gin.H{
					"id":   newTag.Id,
					"name": newTag.Name,
				})
			}
		}
	}
	insertStmt.Close()
	problem.Tags = tags
}

func (problem *Problem) FetchTags() {
	tags := []map[string]interface{}{}
	rows, err := DB.Queryx(`select tag.* from problem_tag inner join tag 
	on problem_tag.tag_id = tag.id where problem_tag.problem_id = ? order by problem_tag.id`, problem.Id)
	if err != nil {
		utils.Consolelog(err)
		return
	}
	for rows.Next() {
		var tag Tag
		err = rows.StructScan(&tag)
		if err != nil {
			utils.Consolelog(err)
		}
		tags = append(tags, gin.H{
			"id":   tag.Id,
			"name": tag.Name,
		})
	}
	rows.Close()
	problem.Tags = tags
}

func (problem *Problem) RemoveTags() error {
	_, err := DB.Exec(`delete from problem_tag where problem_id = ?`, problem.Id)
	problem.Tags = []map[string]interface{}{}
	return err
}

func (problem *Problem) ConvertImgUrl() {
	// 需要将图片地址转换为绝对地址
	if problem.Description.Valid {
		problem.Description.String = utils.ConvertTextImgUrl(problem.Description.String)
	}
	if problem.Input.Valid {
		problem.Input.String = utils.ConvertTextImgUrl(problem.Input.String)
	}
	if problem.Output.Valid {
		problem.Output.String = utils.ConvertTextImgUrl(problem.Output.String)
	}
	if problem.Hint.Valid {
		problem.Hint.String = utils.ConvertTextImgUrl(problem.Hint.String)
	}
}
