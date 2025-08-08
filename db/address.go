package db

import (
	"database/sql"
	"fmt"

	wl_db "github.com/wsva/lib_go_db"
)

func QueryAddress(db *wl_db.DB) ([][]any, error) {
	var rows *sql.Rows
	var err error
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		sqltext := "select b.project_name, a.note, a.address " +
			"from res_passwd_address a, res_project b " +
			"where a.project_id = b.project_id " +
			"order by 1"
		rows, err = db.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", db.Type)
	}
	if err != nil {
		return nil, err
	}
	var result [][]any
	for rows.Next() {
		var f1, f2, f3 sql.NullString
		err = rows.Scan(&f1, &f2, &f3)
		if err != nil {
			return nil, err
		}
		result = append(result, []any{f1.String, f2.String, f3.String})
	}

	return result, rows.Close()
}
