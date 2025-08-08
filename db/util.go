package db

import (
	"database/sql"
	"regexp"

	wl_db "github.com/wsva/lib_go_db"
)

func cleanNote(note string) string {
	regStart := regexp.MustCompile(`^[ ,]+`)
	regEnd := regexp.MustCompile(`[ ,]+$`)
	note = regStart.ReplaceAllString(note, "")
	note = regEnd.ReplaceAllString(note, "")
	return note
}

// queryCount: select count(1) from ...
func count(db *wl_db.DB, queryCount string, args ...any) (int64, error) {
	row, err := db.QueryRow(queryCount, args...)
	if err != nil {
		return 0, err
	}
	var f1 sql.NullInt64
	err = row.Scan(&f1)
	if err != nil {
		return 0, err
	}
	return f1.Int64, nil
}

func getList(db *wl_db.DB, query string, args ...any) ([][]any, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columnCount := len(columnTypes)

	scanArgs := make([]any, columnCount)
	for i, v := range columnTypes {
		switch v.DatabaseTypeName() {
		case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
			scanArgs[i] = new(sql.NullString)
		case "BOOL":
			scanArgs[i] = new(sql.NullBool)
		case "INT4":
			scanArgs[i] = new(sql.NullInt64)
		default:
			scanArgs[i] = new(sql.NullString)
		}
	}

	var result [][]any
	lineCount := 0
	for rows.Next() {
		lineCount++

		line := make([]any, columnCount)
		copy(line, scanArgs)
		err = rows.Scan(line...)
		if err != nil {
			return nil, err
		}

		for k, v := range line {
			switch columnTypes[k].DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				line[k] = v.(*sql.NullString).String
			case "BOOL":
				line[k] = v.(*sql.NullBool).Bool
			case "INT4":
				line[k] = v.(*sql.NullInt64).Int64
			default:
				line[k] = v.(*sql.NullString).String
			}
		}

		result = append(result, line)
	}

	return result, rows.Close()
}
