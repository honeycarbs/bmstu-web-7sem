package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	HEARTBEAT_URL = "/api/heartbeat"
)

func Heartbeat(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.WriteHeader(204)
}
