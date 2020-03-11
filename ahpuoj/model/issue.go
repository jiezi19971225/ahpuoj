package model

import (
	"ahpuoj/utils"
	"database/sql"
	"errors"
)

type Issue struct {
	Id        int    `db:"id" json:"id"`
	Title     string `db:"title" json:"title" binding:"required,max=20"`
	ProblemId int    `db:"problem_id" json:"problem_id" binding:"gte=0"`
	UserId    int    `db:"user_id" json:"user_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	IsDeleted int    `db:"is_deleted" json:"is_deleted"`
	// 附加信息
	Username     string         `db:"username" json:"username"`
	Nick         string         `db:"nick" json:"nick"`
	UserAvatar   string         `db:"avatar" json:"avatar"`
	ReplyCount   int            `db:"reply_count" json:"reply_count"`
	ProblemTitle sql.NullString `db:"ptitle" json:"ptitle"`
}

func (issue *Issue) Save() error {
	result, err := DB.Exec(`insert into issue
	(title,problem_id,user_id,created_at,updated_at) 
	values (?,?,?,NOW(),NOW())`, issue.Title, issue.ProblemId, issue.UserId)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	issue.Id = utils.Int64to32(lastInsertId)
	return err
}

func (issue *Issue) ToggleStatus() error {
	result, err := DB.Exec(`update issue set is_deleted = not is_deleted,updated_at = NOW() where id = ?`, issue.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}
