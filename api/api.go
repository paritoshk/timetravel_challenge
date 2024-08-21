package api

import (
	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/service"
)

type API struct {
	records service.RecordService
}

func NewAPI(records service.RecordService) *API {
	return &API{records}
}

func (a *API) CreateRoutes(router *mux.Router) {
	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/records/{id}", a.GetRecordsV1).Methods("GET")
	v1.HandleFunc("/records/{id}", a.PostRecordsV1).Methods("POST")

	v2 := router.PathPrefix("/api/v2").Subrouter()
	v2.HandleFunc("/records/{id}", a.GetRecordsV2).Methods("GET")
	v2.HandleFunc("/records/{id}", a.PostRecordsV2).Methods("POST")
	v2.HandleFunc("/records/{id}/versions", a.GetRecordVersionsV2).Methods("GET")
}
