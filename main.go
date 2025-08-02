package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	err := initGlobals()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := mux.NewRouter()

	router.Handle("/get/passwd",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGetPassword),
		))
	router.Handle("/get/ip",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGetIP),
		))
	router.Handle("/get/address",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGetAddress),
		))
	router.Handle("/get/unencrypted",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGetUnencrypted),
		))
	router.Handle("/get/encrypted",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGetEncrypted),
		))
	router.Handle("/add",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleAdd),
		))
	router.Handle("/modify",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleModify),
		))
	router.Handle("/encrypt/passwd",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleEncryptPassword),
		))
	router.Handle("/decrypt/passwd",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleDecryptPassword),
		))
	router.Handle("/grant/get/code",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGrantGetCode),
		))
	router.Handle("/grant/do",
		negroni.New(
			negroni.HandlerFunc(handleCheckToken),
			negroni.HandlerFunc(handleGrantDo),
		))

	server := negroni.New(negroni.NewRecovery())
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
