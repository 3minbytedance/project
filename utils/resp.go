package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H struct {
	Code int
	Msg  string
	Data interface{}
}

func resp(w http.ResponseWriter, msg string, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	marshal, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	_, err = w.Write(marshal)
	if err != nil {
		fmt.Println(err)
	}
}

func RespOk(w http.ResponseWriter, msg string, data interface{}) {
	resp(w, msg, 0, data)
}

func RespFail(w http.ResponseWriter, msg string) {
	resp(w, msg, -1, nil)
}
