package main

import (
	"wentmin/components"
	"wentmin/netmodel"
	"wentmin/weblogic"
)

func main() {
	components.InitWebCfg()
	weblogic.RegWSHandlers()
	weblogic.RegHttpHandlers()
	wbs := netmodel.NewWebServer()
	wbs.Start()
}
