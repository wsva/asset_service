package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	wl_int "github.com/wsva/lib_go_integration"
)

func main() {
	err := initGlobals()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := mux.NewRouter()

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir(filepath.Join(Basepath, "template/css/")))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir(filepath.Join(Basepath, "template/js/")))))

	router.Handle("/asset/password",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handlePassword),
		))
	router.Handle("/asset/ip",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleIP),
		))
	router.Handle("/asset/address",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleAddress),
		))
	router.Handle("/permission/code",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGrantGetCode),
		))
	router.Handle("/permission/grant",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGrantDo),
		))

	router.Handle("/",
		negroni.New(
			negroni.HandlerFunc(handleDashboard),
		))
	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(handleLogout),
		))
	router.Handle(wl_int.OAuth2LoginPath,
		negroni.New(
			negroni.HandlerFunc(handleOAuth2Login),
		))
	router.Handle(wl_int.OAuth2CallbackPath,
		negroni.New(
			negroni.HandlerFunc(handleOAuth2Callback),
		))

	server := negroni.New(negroni.NewRecovery())
	server.Use(negroni.NewLogger())
	server.UseHandler(router)

	for _, v := range mainConfig.ListenList {
		if !v.Enable {
			continue
		}
		v1 := v
		switch v1.LowercaseProtocol() {
		case "http":
			go func() {
				err = http.ListenAndServe(fmt.Sprintf(":%v", v1.Port),
					server)
				if err != nil {
					fmt.Println(err)
				}
			}()
		case "https":
			go func() {
				s := &http.Server{
					Addr:    fmt.Sprintf(":%v", v1.Port),
					Handler: server,
				}
				s.SetKeepAlivesEnabled(false)
				err = s.ListenAndServeTLS(ServerCrtFile, ServerKeyFile)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
	}
	select {}
}
