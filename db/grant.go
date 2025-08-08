package db

import (
	"database/sql"
	"errors"
	"fmt"

	wl_db "github.com/wsva/lib_go_db"
)

type GrantCode struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}

type Grant struct {
	//grant, revoke
	Action string `json:"action"`
	//team, project
	Type string `json:"type"`

	AccountID string `json:"account_id"`
	TeamID    string `json:"team_id"`
	ProjectID string `json:"project_id"`
	ClassID   string `json:"class_id"`
}

func (g *Grant) grantTeam(db *wl_db.DB) error {
	switch db.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("insert into was.res_grant_team values "+
			"(sysy_guid(), '%v', '%v', %v)",
			g.AccountID, g.TeamID, g.ClassID)
		_, err := db.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func (g *Grant) grantProject(db *wl_db.DB) error {
	switch db.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("insert into was.res_grant_project values "+
			"(sysy_guid(), '%v', '%v', %v)",
			g.AccountID, g.ProjectID, g.ClassID)
		_, err := db.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func (g *Grant) revokeTeam(db *wl_db.DB) error {
	switch db.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("delete from was.res_grant_team "+
			"where account_id='%v' and team_id='%v'",
			g.AccountID, g.TeamID)
		_, err := db.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func (g *Grant) revokeProject(db *wl_db.DB) error {
	switch db.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("delete from was.res_grant_project "+
			"where account_id='%v' and project_id='%v'",
			g.AccountID, g.ProjectID)
		_, err := db.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func (g *Grant) Do(db *wl_db.DB) error {
	switch g.Action {
	case "grant":
		switch g.Type {
		case "team":
			return g.grantTeam(db)
		case "project":
			return g.grantProject(db)
		default:
			return errors.New("not supported type: " + g.Type)
		}
	case "revoke":
		switch g.Type {
		case "team":
			return g.revokeTeam(db)
		case "project":
			return g.revokeProject(db)
		default:
			return errors.New("not supported type: " + g.Type)
		}
	default:
		return errors.New("unknown action type: " + g.Action)
	}
}

func QueryAccountCode(db *wl_db.DB) ([]*GrantCode, error) {
	switch db.Type {
	case wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := "select account_id, realname note from sys_account where valid='Y'"
		return queryCode(db, query)
	default:
		return nil, fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func QueryProjectCode(db *wl_db.DB) ([]*GrantCode, error) {
	switch db.Type {
	case wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := "select project_id, project_name from res_project"
		return queryCode(db, query)
	default:
		return nil, fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func QueryTeamCode(db *wl_db.DB) ([]*GrantCode, error) {
	switch db.Type {
	case wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := "select team_id, team_name from res_team"
		return queryCode(db, query)
	default:
		return nil, fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func queryCode(db *wl_db.DB, query string, args ...any) ([]*GrantCode, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var result []*GrantCode
	for rows.Next() {
		var f1, f2 sql.NullString
		err = rows.Scan(&f1, &f2)
		if err != nil {
			return nil, err
		}
		result = append(result, &GrantCode{
			ID:   f1.String,
			Note: f2.String,
		})
	}
	rows.Close()
	return result, nil
}
