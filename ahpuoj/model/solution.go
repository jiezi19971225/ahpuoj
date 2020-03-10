package model

import (
	"ahpuoj/utils"
)

type Solution struct {
	Id         int            `db:"solution_id" json:"solution_id"`
	ProblemId  int            `db:"problem_id" json:"problem_id"`
	TeamId     int            `db:"team_id" json:"team_id"`
	UserId     int            `db:"user_id" json:"user_id"`
	ContestId  int            `db:"contest_id" json:"contest_id"`
	Num        int            `db:"num" json:"num"`
	Time       int            `db:"time" json:"time"`
	Memory     int            `db:"memory" json:"memory"`
	InDate     string         `db:"in_date" json:"in_date"`
	Result     int            `db:"result" json:"result"`
	Language   int            `db:"language" json:"language"`
	IP         string         `db:"ip" json:"ip"`
	JudgeTime  NullString 	`db:"judgetime" json:"judgetime"`
	Valid      int            `db:"valid" json:"valid"`
	CodeLength int            `db:"code_length" json:"code_length"`
	PassRate   float32        `db:"pass_rate" json:"pass_rate"`
	LintError  int            `db:"lint_error" json:"lint_error"`
	Judger     NullString 	`db:"judger" json:"judger"`
	// 附加信息
	Public       int    `db:"public" json:"public"`
	Username     string `db:"username" json:"username"`
	Nick         string `db:"nick" json:"nick"`
	UserAvatar   string `db:"avatar" json:"avatar"`
	ProblemTitle string `db:"title" json:"problem_title"`
}

func (solution *Solution) Save() error {
	result, err := DB.Exec(`insert into solution
	(problem_id,team_id,user_id,contest_id,num,in_date,language,ip,code_length,result) 
	values (?,?,?,?,?,NOW(),?,?,?,0)`, solution.ProblemId, solution.TeamId, solution.UserId, solution.ContestId, solution.Num,
		solution.Language, solution.IP, solution.CodeLength)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	solution.Id = utils.Int64to32(lastInsertId)
	return err
}

func (solution *Solution) Response() map[string]interface{} {

	return map[string]interface{}{
		"id":          solution.Id,
		"team_id":     solution.TeamId,
		"contest_id":  solution.ContestId,
		"num":         solution.Num,
		"time":        solution.Time,
		"memory":      solution.Memory,
		"result":      solution.Result,
		"language":    solution.Language,
		"in_date":     solution.InDate,
		"judgetime":   solution.JudgeTime.String,
		"code_length": solution.CodeLength,
		"pass_rate":   solution.PassRate,
		"judger":      solution.Judger.String,
		"public":      solution.Public,
		"user": map[string]interface{}{
			"id":       solution.UserId,
			"username": solution.Username,
			"nick":     solution.Nick,
			"avatar":   solution.UserAvatar,
		},
		"problem": map[string]interface{}{
			"id":    solution.ProblemId,
			"title": solution.ProblemTitle,
		},
	}
}
