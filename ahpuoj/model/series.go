package model

import (
	"ahpuoj/utils"
	"errors"
	"strconv"
)

type Series struct {
	Id          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name" binding:"required,max=20"`
	Description string `db:"description" json:"description"`
	TeamMode    int    `db:"team_mode" json:"team_mode" binding:"gte=0,lte=1"`
	Defunct     int    `db:"defunct" json:"defunct"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
	IsDeleted   int    `db:"is_deleted" json:"is_deleted"`
	UserId      int    `db:"user_id" json:"user_id"`
	Username    string `db:"username" json:"username"`
	// 附加信息
	ContestInfos []map[string]interface{} `json:"contestinfos"`
}

func (series *Series) Save() error {
	result, err := DB.Exec(`insert into series
	(name,description,team_mode,user_id,created_at,updated_at) 
	values (?,?,?,?,NOW(),NOW())`, series.Name, series.Description, series.TeamMode, series.UserId)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	series.Id = utils.Int64to32(lastInsertId)
	return err
}

func (series *Series) Update() error {
	result, err := DB.Exec(`update series set name = ?,description=?,team_mode=?,updated_at = NOW() where id = ?`, series.Name, series.Description, series.TeamMode, series.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (series *Series) Delete() error {
	// 软删除
	result, err := DB.Exec(`update series set  is_deleted = 1 where id = ?`, series.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (series *Series) ToggleStatus() error {
	result, err := DB.Exec(`update series set defunct = not defunct,updated_at = NOW() where id = ?`, series.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (series *Series) AttachContestInfo() {
	contestInfos := make([]map[string]interface{}, 0)
	rows, err := DB.Queryx("select contest.* from contest inner join contest_series on contest_series.contest_id = contest.id where contest_series.series_id = ? and contest.team_mode = ? and contest.is_deleted = 0 and contest.defunct = 0", series.Id, series.TeamMode)
	if err != nil {
		utils.Consolelog(err)
		return
	}
	index := 1
	for rows.Next() {
		var contest Contest
		rows.StructScan(&contest)
		contest.CalcStatus()
		contestInfo := map[string]interface{}{
			"id":     strconv.Itoa(contest.Id),
			"name":   contest.Name,
			"status": contest.Status,
		}
		index++
		contestInfos = append(contestInfos, contestInfo)
	}
	series.ContestInfos = contestInfos
}
