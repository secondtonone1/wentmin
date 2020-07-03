package netmodel

import (
	"fmt"
	"net/http"
	"strconv"
	"wentmin/common"
	"wentmin/components"
)

type WebServer struct {
	webserver *http.Server
}

func (wb *WebServer) ListenAndServe() error {
	address := "0.0.0.0:" + strconv.Itoa(components.WebPort)
	wb.webserver = &http.Server{Addr: address}
	err := wb.webserver.ListenAndServe()
	return err
}

func (wb *WebServer) Start() {
	fmt.Println("web server start...")
	err := wb.ListenAndServe()
	if err != nil {
		fmt.Println(common.ErrWebListenFailed.Error())
		return
	}
}

func (wb *WebServer) Stop() {
	wb.webserver.Shutdown(nil)
}

func NewWebServer() *WebServer {
	return &WebServer{}
}
