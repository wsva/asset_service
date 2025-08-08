package db

import (
	"database/sql"
	"fmt"

	wl_db "github.com/wsva/lib_go_db"
)

type IPNote struct {
	IP   string `json:"ip"`
	Note string `json:"note"`
}

func QueryIP(db *wl_db.DB) ([]IPNote, error) {
	var rows *sql.Rows
	var err error
	var result []IPNote
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "select ip, note from v_res_all_ip"
		rows, err = db.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", db.Type)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f1, f2 sql.NullString
		err = rows.Scan(&f1, &f2)
		if err != nil {
			return nil, err
		}
		res := IPNote{
			IP:   f1.String,
			Note: cleanNote(f2.String),
		}
		result = append(result, res)
	}
	rows.Close()
	return result, nil
}
