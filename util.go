package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	wl_db "github.com/wsva/lib_go_db"
	wl_int "github.com/wsva/lib_go_integration"
)

func cleanNote(note string) string {
	regStart := regexp.MustCompile(`^[ ,]+`)
	regEnd := regexp.MustCompile(`[ ,]+$`)
	note = regStart.ReplaceAllString(note, "")
	note = regEnd.ReplaceAllString(note, "")
	return note
}

func sqlsafe(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func getPasswdFromDatabase(r *http.Request) ([]Passwd, error) {
	token, err := wl_int.ParseTokenFromRequest(r)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	var result []Passwd
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("select R.PROJECT_NAME, R.ADDRESS, R.NOTE, "+
			"R.USERNAME, R.PASSWD, R.UUID, R.ENCRYPTED "+
			"from gm.v_res_All_Granted_Passwd R, gm.sys_Token T "+
			"where R.Account_ID=T.Account_ID and T.Token='%v'",
			sqlsafe(token))
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f1, f2, f3, f4, f5, f6, f7 sql.NullString
		err = rows.Scan(&f1, &f2, &f3, &f4, &f5, &f6, &f7)
		if err != nil {
			return nil, err
		}
		res := Passwd{
			Project:   f1.String,
			Address:   f2.String,
			Note:      f3.String,
			Username:  f4.String,
			Passwd:    f5.String,
			UUID:      f6.String,
			Encrypted: f7.String,
		}
		if err := res.Process(); err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	rows.Close()
	return result, nil
}

func getIPFromDatabase() ([]IPNote, error) {
	var rows *sql.Rows
	var err error
	var result []IPNote
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "select A.IP, A.Note from gm.v_res_All_IP A"
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
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
			Note: f2.String,
		}
		res.Process()
		result = append(result, res)
	}
	rows.Close()
	return result, nil
}

func getUnencryptedCountFromDatabase() (int64, error) {
	var rows *sql.Rows
	var err error
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "select count(1) from gm.v_res_All_Plain_Passwd"
		rows, err = cc.DB.Query(sqltext)
	default:
		return 0, fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
	if err != nil {
		return 0, err
	}
	var f1 sql.NullInt64
	rowsCount := 0
	for rows.Next() {
		err = rows.Scan(&f1)
		if err != nil {
			return 0, err
		}
		rowsCount++
	}
	if rowsCount < 1 {
		return 0, errors.New("no data found")
	}
	if rowsCount > 1 {
		return 0, errors.New("conflict: more than one rows found")
	}
	if err != nil {
		return 0, err
	}
	return f1.Int64, nil
}

func getEncryptedCountFromDatabase() (int64, error) {
	var rows *sql.Rows
	var err error
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "select count(1) from gm.v_res_All_Encrypted_Passwd"
		rows, err = cc.DB.Query(sqltext)
	default:
		return 0, fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
	if err != nil {
		return 0, err
	}
	var f1 sql.NullInt64
	rowsCount := 0
	for rows.Next() {
		err = rows.Scan(&f1)
		if err != nil {
			return 0, err
		}
		rowsCount++
	}
	if rowsCount < 1 {
		return 0, errors.New("no data found")
	}
	if rowsCount > 1 {
		return 0, errors.New("conflict: more than one rows found")
	}
	if err != nil {
		return 0, err
	}
	return f1.Int64, nil
}

func getPlainPasswdListFromDatabase() ([]*Passwd4Update, error) {
	var rows *sql.Rows
	var err error
	var result []*Passwd4Update
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "select Table_Name, UUID, Passwd " +
			"from gm.v_res_All_Plain_Passwd"
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
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
		res := Passwd4Update{
			TableName: f1.String,
			UUID:      f2.String,
			Passwd:    f3.String,
		}
		result = append(result, &res)
	}
	rows.Close()
	return result, nil
}

func getEncryptedPasswdListFromDatabase() ([]*Passwd4Update, error) {
	var rows *sql.Rows
	var err error
	var result []*Passwd4Update
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "select Table_Name, UUID, Passwd " +
			"from gm.v_res_All_Encrypted_Passwd"
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
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
		res := Passwd4Update{
			TableName: f1.String,
			UUID:      f2.String,
			Passwd:    f3.String,
		}
		result = append(result, &res)
	}
	rows.Close()
	return result, nil
}

func getAddressListFromDatabase() ([]ProjectAddress, error) {
	var rows *sql.Rows
	var err error
	resultMap := make(map[string][]Address)
	var resultList []ProjectAddress
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle:
		sqltext := "SELECT b.project_name project, a.Note Note, a.Address " +
			"FROM gm.res_Passwd_Address a, gm.res_Project b " +
			"WHERE a.Project_ID = b.Project_ID " +
			"order by 1"
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
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

func getAccountCodeListFromDatabase() ([]GrantCode, error) {
	var rows *sql.Rows
	var err error
	var result []GrantCode
	switch cc.DB.Type {
	case wl_db.DBTypeOracle:
		sqltext := "SELECT account_id, " +
			"realname||', '||phonenumber||', Valid:'||valid note " +
			"FROM gm.sys_account"
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
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
		result = append(result, GrantCode{
			ID:   f1.String,
			Note: f2.String,
		})
	}
	rows.Close()
	return result, nil
}

func getProjectCodeListFromDatabase() ([]GrantCode, error) {
	var rows *sql.Rows
	var err error
	var result []GrantCode
	switch cc.DB.Type {
	case wl_db.DBTypeOracle:
		sqltext := "SELECT project_id, project_name " +
			"FROM gm.res_project"
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
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
		result = append(result, GrantCode{
			ID:   f1.String,
			Note: f2.String,
		})
	}
	rows.Close()
	return result, nil
}

func getTeamCodeListFromDatabase() ([]GrantCode, error) {
	var rows *sql.Rows
	var err error
	var result []GrantCode
	switch cc.DB.Type {
	case wl_db.DBTypeOracle:
		sqltext := "SELECT team_id, team_name FROM gm.res_team"
		rows, err = cc.DB.Query(sqltext)
	default:
		return nil, fmt.Errorf("invalid DBType %v", cc.DB.Type)
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
		result = append(result, GrantCode{
			ID:   f1.String,
			Note: f2.String,
		})
	}
	rows.Close()
	return result, nil
}
