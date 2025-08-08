package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	wl_http "github.com/wsva/lib_go/http"
	wl_int "github.com/wsva/lib_go_integration"

	"github.com/wsva/asset_service/db"
)

func handlePassword(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	action := r.FormValue("action")
	account_id := r.FormValue("account_id")
	switch action {
	case "project":
		granted := db.NewGranted(&cc.DB, account_id, AESKey)
		list, err := granted.ProjectList()
		if err != nil {
			fmt.Println(err)
			wl_http.RespondError(w, "database error")
			return
		}
		fmt.Println(list)
		wl_http.RespondJSON(w, wl_http.Response{
			Success: true,
			Data: wl_http.ResponseData{
				List: list,
			},
		})
	case "description":
		project_id := r.FormValue("project_id")
		granted := db.NewGranted(&cc.DB, account_id, AESKey)
		list, err := granted.DescriptionList(project_id)
		if err != nil {
			fmt.Println(err)
			wl_http.RespondError(w, "database error")
			return
		}
		wl_http.RespondJSON(w, wl_http.Response{
			Success: true,
			Data: wl_http.ResponseData{
				List: list,
			},
		})
	case "password":
		uuid := r.FormValue("uuid")
		granted := db.NewGranted(&cc.DB, account_id, AESKey)
		passwd, err := granted.Password(uuid)
		if err != nil {
			fmt.Println(err)
			wl_http.RespondError(w, "database error")
			return
		}
		wl_http.RespondJSON(w, wl_http.Response{
			Success: true,
			Data: wl_http.ResponseData{
				List: []string{passwd},
			},
		})
	case "encrypt":
		if encryptLock.TryLock() {
			defer encryptLock.Unlock()
			err := db.EncryptPassword(&cc.DB, AESKey)
			if err != nil {
				fmt.Println(err)
				wl_http.RespondError(w, "database error")
				return
			}
			wl_http.RespondSuccess(w)
		} else {
			wl_http.RespondError(w, "action locked")
		}
	case "decrypt":
		if encryptLock.TryLock() {
			defer encryptLock.Unlock()
			err := db.DecryptPassword(&cc.DB, AESKey)
			if err != nil {
				fmt.Println(err)
				wl_http.RespondError(w, "database error")
				return
			}
			wl_http.RespondSuccess(w)
		} else {
			wl_http.RespondError(w, "action locked")
		}
	case "overview":
		encrypted, err := db.CountEncryptedPassword(&cc.DB)
		if err != nil {
			fmt.Println(err)
			wl_http.RespondError(w, "database error")
			return
		}
		unencrypted, err := db.CountUnencryptedPassword(&cc.DB)
		if err != nil {
			fmt.Println(err)
			wl_http.RespondError(w, "database error")
			return
		}
		wl_http.RespondJSON(w, wl_http.Response{
			Success: true,
			Data: wl_http.ResponseData{
				List: []map[string]any{
					{"name": "encrypted", "value": encrypted},
					{"name": "unencrypted", "value": unencrypted},
				},
			},
		})
	default:
		wl_http.RespondError(w, "invalid action")
	}
}

func handleIP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resList, err := db.QueryIP(&cc.DB)
	if err != nil {
		fmt.Println(err)
		wl_http.RespondError(w, "database error")
		return
	}
	wl_http.RespondJSON(w, wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: resList,
		},
	})
}

func handleGetAddress(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resList, err := db.QueryAddress(&cc.DB)
	if err != nil {
		fmt.Println(err)
		wl_http.RespondError(w, "database error")
		return
	}
	wl_http.RespondJSON(w, wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: resList,
		},
	})
}

func handleGrantGetCode(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accountList, err := db.QueryAccountCode(&cc.DB)
	if err != nil {
		fmt.Println(err)
		wl_http.RespondError(w, "database error")
		return
	}
	projectList, err := db.QueryProjectCode(&cc.DB)
	if err != nil {
		fmt.Println(err)
		wl_http.RespondError(w, "database error")
		return
	}
	teamList, err := db.QueryTeamCode(&cc.DB)
	if err != nil {
		fmt.Println(err)
		wl_http.RespondError(w, "database error")
		return
	}
	wl_http.RespondJSON(w, wl_http.Response{
		Success: true,
		Data: wl_http.ResponseData{
			List: []any{
				accountList,
				projectList,
				teamList,
			},
		},
	})
}

func handleGrantDo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	req, err := wl_http.ParseRequest(r, 1024)
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	var grant db.Grant
	err = json.Unmarshal(req.Data, &grant)
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	err = grant.Do(&cc.DB)
	if err != nil {
		wl_http.RespondError(w, err)
		return
	}
	wl_http.RespondSuccess(w)
}

func handleOAuth2Login(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	oa := oaMap.Add(&mainConfig.AuthService, httpsClient)
	oa.HandleLogin(w, r)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	state := r.FormValue("state")
	oa, err := oaMap.Get(state)
	if err != nil {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	oaMap.Delete(state)
	oa.HandleCallback(w, r)
}

func handleDashboard(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tpFile := filepath.Join(Basepath, "template/html/dashboard.html")
	tp, err := template.ParseFiles(tpFile)
	if err != nil {
		fmt.Fprintf(w, "parse template %v error: %v", tpFile, err)
		return
	}

	type Data struct {
		Name  string
		Email string
	}

	if wl_int.VerifyToken(r, httpsClient, mainConfig.AuthService.IntrospectURL) != nil {
		tp.Execute(w, Data{})
		return
	}

	userinfoCookie, err := r.Cookie("userinfo")
	if err != nil {
		wl_http.RespondError(w, "missing user info")
		return
	}
	userinfoBytes, err := base64.URLEncoding.DecodeString(userinfoCookie.Value)
	if err != nil {
		wl_http.RespondError(w, "invalid user info")
		return
	}

	var userinfo wl_int.UserInfo
	err = json.Unmarshal(userinfoBytes, &userinfo)
	if err != nil {
		wl_http.RespondError(w, "invalid user info")
		return
	}

	tp.Execute(w, Data{
		Name:  userinfo.Name,
		Email: userinfo.Email,
	})
}

func handleLogout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	wl_int.DeleteCookieToken(w, "access_token")
	wl_int.DeleteCookieToken(w, "refresh_token")
	wl_int.DeleteCookieToken(w, "userinfo")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleCheckToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !wl_int.CheckInternalKey(r, AESKey, AESIV) {
		if err := wl_int.VerifyToken(r, httpsClient, mainConfig.AuthService.IntrospectURL); err != nil {
			fmt.Println("verify token error: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	next(w, r)
}
