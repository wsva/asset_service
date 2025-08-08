package db

import (
	"database/sql"
	"fmt"

	wl_db "github.com/wsva/lib_go_db"
)

func QueryIP(db *wl_db.DB) ([][]any, error) {
	var rows *sql.Rows
	var err error
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		sqltext := "select ip, note from v_res_all_ip"
		rows, err = db.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", db.Type)
	}
	if err != nil {
		return nil, err
	}
	var result [][]any
	for rows.Next() {
		var f1, f2 sql.NullString
		err = rows.Scan(&f1, &f2)
		if err != nil {
			return nil, err
		}
		result = append(result, []any{f1.String, cleanNote(f2.String)})
	}
	return result, rows.Close()
}
