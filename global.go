package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"sync"

	wl_fs "github.com/wsva/lib_go/fs"
	wl_http "github.com/wsva/lib_go/http"
	wl_int "github.com/wsva/lib_go_integration"
)

const (
	AESKey = "key"
	AESIV  = "iv"
)

type MainConfig struct {
	ListenList  []wl_http.ListenInfo `json:"ListenList"`
	AuthService wl_int.AuthService   `json:"AuthService"`
}

var (
	Basepath       = ""
	MainConfigFile = path.Join(wl_int.DirConfig, "asset_service_config.json")
	CACrtFile      = path.Join(wl_int.DirPKI, wl_int.CACrtFile)
	ServerCrtFile  = path.Join(wl_int.DirPKI, wl_int.ServerCrtFile)
	ServerKeyFile  = path.Join(wl_int.DirPKI, wl_int.ServerKeyFile)
)

var mainConfig MainConfig
var cc *wl_int.CommonConfig
var httpsClient *http.Client
var oaMap wl_int.OAuth2Map
var encryptLock sync.Mutex

func initGlobals() error {
	basepath, err := wl_fs.GetExecutableFullpath()
	if err != nil {
		return err
	}
	Basepath = basepath
	MainConfigFile = path.Join(basepath, MainConfigFile)

	CACrtFile = path.Join(basepath, CACrtFile)
	ServerCrtFile = path.Join(basepath, ServerCrtFile)
	ServerKeyFile = path.Join(basepath, ServerKeyFile)

	contentBytes, err := os.ReadFile(MainConfigFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contentBytes, &mainConfig)
	if err != nil {
		return err
	}
	cc, err = wl_int.LoadCommonConfig(basepath, AESKey, AESIV)
	if err != nil {
		return err
	}

	httpsClient, err = wl_int.InitHttpsClient(CACrtFile)
	if err != nil {
		return err
	}

	return nil
}
