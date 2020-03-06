package model

import (
	"ahpuoj/utils"
	"encoding/json"
	"errors"
)

type Reply struct {
	Id          int    `db:"id" json:"id"`
	IssueId     int    `db:"issue_id" json:"issue_id"`
	UserId      int    `db:"user_id" json:"user_id"`
	ReplyId     int    `db:"reply_id" json:"reply_id" binding:"gte=0"`
	ReplyUserId int    `db:"reply_user_id" json:"reply_user_id" binding:"gte=0"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
	Content     string `db:"content" json:"content" binding:"required"`
	IsDeleted   int    `db:"is_deleted" json:"is_deleted"`
	Status      int    `db:"status" json:"status"`
	// 附加信息
	Username      string  `db:"username" json:"username"`
	ReplyUserNick string  `db:"rnick" json:"rnick"`
	Nick          string  `db:"nick" json:"nick"`
	UserAvatar    string  `db:"avatar" json:"avatar"`
	ReplyCount    int     `db:"reply_count" json:"reply_count"`
	IssueTitle    string  `db:"issue_title" json:"issue_title"`
	SubReplys     []Reply `json:"sub_replys"`
}

func (reply *Reply) MarshalJSON() ([]byte, error) {
	type Alias Reply
	reply.Content = utils.ConvertTextImgUrl(reply.Content)
	return json.Marshal((*Alias)(reply))
}

func (reply *Reply) Save() error {
	result, err := DB.Exec(`insert into reply
	(user_id,issue_id,reply_id,reply_user_id,content,created_at,updated_at) 
	values (?,?,?,?,?,NOW(),NOW())`, reply.UserId, reply.IssueId, reply.ReplyId, reply.ReplyUserId, reply.Content)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	reply.Id = utils.Int64to32(lastInsertId)
	// 更新主题的最后更新时间
	DB.Exec("update issue set updated_at = NOW() where id = ?", reply.IssueId)
	return err
}

func (reply *Reply) ToggleStatus() error {
	result, err := DB.Exec(`update reply set is_deleted = not is_deleted,updated_at = NOW() where id = ?`, reply.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}
