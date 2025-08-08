package db

import (
	"database/sql"
	"fmt"
	"sort"

	wl_db "github.com/wsva/lib_go_db"
)

type Address struct {
	Note    string `json:"note"`
	Address string `json:"address"`
}

type ProjectAddress struct {
	Project     string    `json:"project"`
	AddressList []Address `json:"list"`
}

func QueryAddress(db *wl_db.DB) ([]ProjectAddress, error) {
	var rows *sql.Rows
	var err error
	resultMap := make(map[string][]Address)
	var resultList []ProjectAddress
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "SELECT b.project_name project, a.Note Note, a.Address " +
			"FROM gm.res_Passwd_Address a, gm.res_Project b " +
			"WHERE a.Project_ID = b.Project_ID " +
			"order by 1"
		rows, err = db.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", db.Type)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f1, f2, f3 sql.NullString
		err = rows.Scan(&f1, &f2, &f3)
		if err != nil {
			return nil, err
		}
		resultMap[f1.String] = append(resultMap[f1.String], Address{
			Note:    f2.String,
			Address: f3.String,
		})
	}
	rows.Close()
	for k, v := range resultMap {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Address < v[j].Address
		})
		resultList = append(resultList, ProjectAddress{
			Project:     k,
			AddressList: resultMap[k],
		})
	}
	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].Project < resultList[j].Project
	})
	return resultList, nil
}
