package model

import (
	"ahpuoj/service/mysql"
	"ahpuoj/utils"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func init() {
	DB = mysql.Pool
}

type NullString struct {
	sql.NullString
}

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}
func (v *NullString) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))

	if text == "null" {
		v.Valid = false
		return nil
	}

	err := json.Unmarshal(b, &v.String)
	if err != nil {
		return err
	}

	v.Valid = true
	return nil
}

func Paginate(page *int, perpage *int, tableName string, outputField []string, whereString string) (*sqlx.Rows, int, error) {
	if *page < 0 {
		*page = 1
	}
	if *perpage < 0 {
		*perpage = 0
	}
	var total int
	sql := "select count(1) from  (select 1 from " + tableName + " " + whereString + ")pesudo"
	err := DB.Get(&total, sql)
	var maxPage int
	if total == 0 {
		maxPage = 1
	} else if total%*perpage == 0 {
		maxPage = total / *perpage
	} else {
		maxPage = total / *perpage + 1
	}
	if *page > maxPage {
		*page = maxPage
	}
	var rows *sqlx.Rows
	offset := (*page - 1) * *perpage
	offsetStr := strconv.Itoa(offset)
	limitStr := strconv.Itoa(*perpage)
	utils.Consolelog(err, total, maxPage, *page, *perpage)
	outputFieldString := strings.Join(outputField, ",")
	sql = "select " + outputFieldString + " from " + tableName + " " + whereString + " limit " + limitStr + " offset " + offsetStr
	utils.Consolelog(sql)
	rows, err = DB.Queryx(sql)
	if err != nil {
		utils.Consolelog(err)
	}
	return rows, total, err
}
