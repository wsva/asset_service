package main

import (
	"fmt"

	wl_crypto "github.com/wsva/lib_go/crypto"
	wl_db "github.com/wsva/lib_go_db"
)

type Passwd struct {
	Project   string `json:"project"`
	Address   string `json:"address"`
	Note      string `json:"note"`
	Username  string `json:"username"`
	Passwd    string `json:"passwd"`
	UUID      string `json:"uuid"`
	Encrypted string `json:"-"`
}

func (r *Passwd) Process() error {
	r.Note = cleanNote(r.Note)

	if r.Encrypted == "Y" {
		plainText, err := wl_crypto.AES256Decrypt(AESKey, r.UUID, r.Passwd)
		if err != nil {
			return err
		}
		r.Passwd = plainText
	}

	return nil
}

type Passwd4Update struct {
	TableName string
	UUID      string
	Passwd    string
}

func (p *Passwd4Update) Encrypt() error {
	ciphertext, err := wl_crypto.AES256Encrypt(AESKey, p.UUID, p.Passwd)
	if err != nil {
		return err
	}
	sqltext := fmt.Sprintf("update %v "+
		"set Passwd='%v', Encrypted='Y' where UUID='%v'",
		p.TableName, ciphertext, p.UUID)
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL:
		_, err = cc.DB.Exec(sqltext)
	case wl_db.DBTypeOracle:
		_, err = cc.DB.Exec(sqltext)
	default:
		return fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
	return err
}

func (p *Passwd4Update) Decrypt() error {
	plaintext, err := wl_crypto.AES256Decrypt(AESKey, p.UUID, p.Passwd)
	if err != nil {
		return err
	}
	sqltext := fmt.Sprintf("update %v "+
		"set Passwd='%v', Encrypted='N' where UUID='%v'",
		p.TableName, plaintext, p.UUID)
	switch cc.DB.Type {
	case wl_db.DBTypeMySQL:
		_, err = cc.DB.Exec(sqltext)
	case wl_db.DBTypeOracle:
		_, err = cc.DB.Exec(sqltext)
	default:
		return fmt.Errorf("invalid DBType %v", cc.DB.Type)
	}
	return err
}

type Address struct {
	Note    string `json:"note"`
	Address string `json:"address"`
}

type ProjectAddress struct {
	Project     string    `json:"project"`
	AddressList []Address `json:"list"`
}
