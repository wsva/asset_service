package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	wl_http "github.com/wsva/lib_go/http"
	wl_int "github.com/wsva/lib_go_integration"
)

func handleGetPassword(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resList, err := getPasswdFromDatabase(r)
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	resp := wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: resList,
		},
	}
	resp.DoResponse(w)
}

func handleGetIP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resList, err := getIPFromDatabase()
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	resp := wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: resList,
		},
	}
	resp.DoResponse(w)
}

func handleGetAddress(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resList, err := getAddressListFromDatabase()
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	resp := wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: resList,
		},
	}
	resp.DoResponse(w)
}

func handleGetUnencrypted(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	result, err := getUnencryptedCountFromDatabase()
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	resp := wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: []int64{result},
		},
	}
	resp.DoResponse(w)
}

func handleGetEncrypted(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	result, err := getEncryptedCountFromDatabase()
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	resp := wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: []int64{result},
		},
	}
	resp.DoResponse(w)
}

func handleAdd(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
}

func handleModify(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
}

func handleEncryptPassword(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resList, err := getPlainPasswdListFromDatabase()
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	for _, v := range resList {
		err = v.Encrypt()
		if err != nil {
			wl_http.RespondError(w, err)
			return
		}
	}
	result, err := getUnencryptedCountFromDatabase()
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	resp := wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: []int64{result},
		},
	}
	resp.DoResponse(w)
}

func handleDecryptPassword(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resList, err := getEncryptedPasswdListFromDatabase()
	if err != nil {
		fmt.Println("query databse error: ", err)
		wl_http.RespondError(w, err)
		return
	}
	for _, v := range resList {
		err = v.Decrypt()
		if err != nil {
			wl_http.RespondError(w, err)
			return
		}
	}
	resp := wl_http.Response{
		Success: true,
	}
	resp.DoResponse(w)
}

func handleGrantGetCode(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accountList, err := getAccountCodeListFromDatabase()
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	projectList, err := getProjectCodeListFromDatabase()
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	teamList, err := getTeamCodeListFromDatabase()
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	resp := wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: []interface{}{
				accountList,
				projectList,
				teamList,
			},
		},
	}
	resp.DoResponse(w)
}

func handleGrantDo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	req, err := wl_http.ParseRequest(r, 1024)
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	var grant Grant
	err = json.Unmarshal(req.Data, &grant)
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	err = grant.Do()
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	wl_http.RespondSuccess(w)
}

func handleCheckToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !wl_int.CheckInternalKey(r, AESKey, AESIV) {
		token, err := wl_int.ParseTokenFromRequest(r)
		if err != nil {
			fmt.Println("parse token error: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = wl_int.CheckAndRefreshToken(cc.AccountAddress, CACrtFile, token)
		if err != nil {
			fmt.Println("check token error: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	next(w, r)
}
