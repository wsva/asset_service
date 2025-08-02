package main

import (
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

func (g *Grant) grantTeam() error {
	switch cc.DB.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("insert into was.res_grant_team values "+
			"(sysy_guid(), '%v', '%v', %v)",
			g.AccountID, g.TeamID, g.ClassID)
		_, err := cc.DB.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
}

func (g *Grant) grantProject() error {
	switch cc.DB.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("insert into was.res_grant_project values "+
			"(sysy_guid(), '%v', '%v', %v)",
			g.AccountID, g.ProjectID, g.ClassID)
		_, err := cc.DB.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
}

func (g *Grant) revokeTeam() error {
	switch cc.DB.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("delete from was.res_grant_team "+
			"where account_id='%v' and team_id='%v'",
			g.AccountID, g.TeamID)
		_, err := cc.DB.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
}

func (g *Grant) revokeProject() error {
	switch cc.DB.Type {
	case wl_db.DBTypeOracle:
		sqltext := fmt.Sprintf("delete from was.res_grant_project "+
			"where account_id='%v' and project_id='%v'",
			g.AccountID, g.ProjectID)
		_, err := cc.DB.Exec(sqltext)
		return err
	default:
		return fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
}

func (g *Grant) Do() error {
	switch g.Action {
	case "grant":
		switch g.Type {
		case "team":
			return g.grantTeam()
		case "project":
			return g.grantProject()
		default:
			return errors.New("not supported type: " + g.Type)
		}
	case "revoke":
		switch g.Type {
		case "team":
			return g.revokeTeam()
		case "project":
			return g.revokeProject()
		default:
			return errors.New("not supported type: " + g.Type)
		}
	default:
		return errors.New("unknown action type: " + g.Action)
	}

}
