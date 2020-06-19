package weblogic

import (
	"net/http"
)

func RegHttpHandlers() {
	http.HandleFunc("/httptest", UsrInfoReq)
}
