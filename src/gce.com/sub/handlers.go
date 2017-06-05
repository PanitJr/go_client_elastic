package sub

import (
	"encoding/json"
	"fmt"
	"go_client_elastic/pkg/go_client_elastic"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

//Index is a
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
	//fmt.Fprintln(w, "data : ", data)
}

//GetLastRecord is a
func GetLastRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	app := vars["app"]
	requestMap := go_client_elastic.AppHostParam{}
	requestMap.Param = append(requestMap.Param, go_client_elastic.AppHost{App: app, Host: host})
	request := go_client_elastic.AppHostReqstBuilder{
		Host:         "10.138.32.97",
		Port:         "8080",
		Index:        "interstellar-crawler",
		Type:         "log",
		AppHostParam: requestMap,
	}
	//json.NewEncoder(w).Encode(request)
	data := go_client_elastic.GetByAppHost(request)
	json.NewEncoder(w).Encode(data)
}

func GetLastRecordSet(w http.ResponseWriter, r *http.Request) {
	var request go_client_elastic.AppHostReqstBuilder
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	data := go_client_elastic.GetByAppHost(request)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
