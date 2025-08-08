package db

import (
	"database/sql"
	"fmt"

	wl_crypto "github.com/wsva/lib_go/crypto"
	wl_db "github.com/wsva/lib_go_db"
)

type Granted struct {
	DB        *wl_db.DB
	AccountID string
	AESkey    string
}

func NewGranted(db *wl_db.DB, account_id, aesKey string) *Granted {
	return &Granted{
		DB:        db,
		AccountID: account_id,
		AESkey:    aesKey,
	}
}

func (p *Granted) ProjectList() ([][]any, error) {
	switch p.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := fmt.Sprintf("select distinct project_id, project_name " +
			"from v_res_all_granted_passwd " +
			"where account_id=$1")
		return getList(p.DB, query, p.AccountID)
	default:
		return nil, fmt.Errorf("invalid DBType %v", p.DB.Type)
	}
}

func (p *Granted) DescriptionList(project_id string) ([][]any, error) {
	switch p.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := fmt.Sprintf("select uuid, " +
			"CONCAT(address, ' - ', note, ' - ', username) AS description " +
			"from v_res_all_granted_passwd " +
			"where project_id=$1 and account_id=$2")
		return getList(p.DB, query, project_id, p.AccountID)
	default:
		return nil, fmt.Errorf("invalid DBType %v", p.DB.Type)
	}
}

func (p *Granted) Password(uuid string) (string, error) {
	var row *sql.Row
	var err error
	switch p.DB.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		sqltext := fmt.Sprintf("select passwd, encrypted " +
			"from v_res_all_granted_passwd " +
			"where uuid=$1")
		row, err = p.DB.QueryRow(sqltext, uuid)
	default:
		return "", fmt.Errorf("invalid DBType %v", p.DB.Type)
	}
	if err != nil {
		return "", err
	}
	var f1, f2 sql.NullString
	err = row.Scan(&f1, &f2)
	if err != nil {
		return "", err
	}
	if f2.String == "Y" {
		return wl_crypto.AES256Decrypt(p.AESkey, uuid, f1.String)
	}
	return f1.String, nil
}

func CountUnencryptedPassword(db *wl_db.DB) (int64, error) {
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := "select count(1) from v_res_all_plain_passwd"
		return count(db, query)
	default:
		return 0, fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func CountEncryptedPassword(db *wl_db.DB) (int64, error) {
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := "select count(1) from v_res_all_encrypted_passwd"
		return count(db, query)
	default:
		return 0, fmt.Errorf("invalid DBType %v", db.Type)
	}
}

func EncryptPassword(db *wl_db.DB, aesKey string) error {
	var rows *sql.Rows
	var err error
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := "select table_name, uuid, passwd from v_res_all_plain_passwd"
		rows, err = db.Query(query)
	default:
		return fmt.Errorf("invalid DBType %v", db.Type)
	}
	if err != nil {
		return err
	}

	encrypt := func(table, uuid, passwd string) error {
		ciphertext, err := wl_crypto.AES256Encrypt(aesKey, uuid, passwd)
		if err != nil {
			return err
		}
		switch db.Type {
		case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
			query := "update $1 set passwd=$2, encrypted='Y' where uuid=$3"
			_, err = db.Exec(query, table, ciphertext, uuid)
		default:
			return fmt.Errorf("invalid DBType %v", db.Type)
		}
		return err
	}

	for rows.Next() {
		var f1, f2, f3 sql.NullString
		if rows.Scan(&f1, &f2, &f3) != nil {
			return err
		}
		if encrypt(f1.String, f2.String, f3.String) != nil {
			return err
		}
	}
	return rows.Close()
}

func DecryptPassword(db *wl_db.DB, aesKey string) error {
	var rows *sql.Rows
	var err error
	switch db.Type {
	case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
		query := "select table_name, uuid, passwd from v_res_all_encrypted_passwd"
		rows, err = db.Query(query)
	default:
		return fmt.Errorf("invalid DBType %v", db.Type)
	}
	if err != nil {
		return err
	}

	decrypt := func(table, uuid, passwd string) error {
		plaintext, err := wl_crypto.AES256Decrypt(aesKey, uuid, passwd)
		if err != nil {
			return err
		}
		switch db.Type {
		case wl_db.DBTypeMySQL, wl_db.DBTypeOracle, wl_db.DBTypePostgreSQL:
			query := "update $1 set passwd=$2, encrypted='N' where uuid=$3"
			_, err = db.Exec(query, table, plaintext, uuid)
		default:
			return fmt.Errorf("invalid DBType %v", db.Type)
		}
		return err
	}

	for rows.Next() {
		var f1, f2, f3 sql.NullString
		if rows.Scan(&f1, &f2, &f3) != nil {
			return err
		}
		if decrypt(f1.String, f2.String, f3.String) != nil {
			return err
		}
	}
	return rows.Close()
}
