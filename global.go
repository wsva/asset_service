package main

import (
	"encoding/json"
	"os"
	"path"

	wl_fs "github.com/wsva/lib_go/fs"
	wl_http "github.com/wsva/lib_go/http"
	wl_int "github.com/wsva/lib_go_integration"
)

const (
	AESKey = "key"
	AESIV  = "iv"
)

type MainConfig struct {
	ListenList []wl_http.ListenInfo `json:"ListenList"`
}

var (
	MainConfigFile = path.Join(wl_int.DirConfig, "assets_management_config.json")
	PKIPath        = wl_int.DirPKI
	CACrtFile      = path.Join(wl_int.DirPKI, wl_int.CACrtFile)
	ServerCrtFile  = path.Join(wl_int.DirPKI, wl_int.ServerCrtFile)
	ServerKeyFile  = path.Join(wl_int.DirPKI, wl_int.ServerKeyFile)
)

var mainConfig MainConfig
var cc *wl_int.CommonConfig

func initGlobals() error {
	basepath, err := wl_fs.GetExecutableFullpath()
	if err != nil {
		return err
	}
	MainConfigFile = path.Join(basepath, MainConfigFile)

	PKIPath = path.Join(basepath, PKIPath)
	CACrtFile = path.Join(PKIPath, CACrtFile)
	ServerCrtFile = path.Join(PKIPath, ServerCrtFile)
	ServerKeyFile = path.Join(PKIPath, ServerKeyFile)

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
	return nil
}
