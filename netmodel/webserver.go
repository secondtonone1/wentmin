package netmodel

import (
	"fmt"
	"net/http"
	"strconv"
	"wentmin/components"
	"wentserver/config"
)

type WebServer struct {
}

func (wb *WebServer) ListenAndServe() error {
	address := "0.0.0.0:" + strconv.Itoa(components.WebPort)
	err := http.ListenAndServe(address, nil)
	return err
}

func (wb *WebServer) Start() {
	fmt.Println("web server start...")
	err := wb.ListenAndServe()
	if err != nil {
		fmt.Println(config.ErrWebListenFailed.Error())
		return
	}
}

func NewWebServer() *WebServer {
	return &WebServer{}
}
