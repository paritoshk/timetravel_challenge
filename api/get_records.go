package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *API) GetRecordsV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)
	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	record, err := a.records.GetRecord(ctx, int(idNumber))
	if err != nil {
		err := writeError(w, fmt.Sprintf("failed to get record: %v", err), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}

func (a *API) GetRecordsV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	version := r.URL.Query().Get("version")

	idNumber, err := strconv.ParseInt(id, 10, 32)
	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	var record interface{}
	var getErr error

	if version == "" {
		record, getErr = a.records.GetRecord(ctx, int(idNumber))
	} else {
		versionNumber, err := strconv.ParseInt(version, 10, 32)
		if err != nil || versionNumber <= 0 {
			err := writeError(w, "invalid version; version must be a positive number", http.StatusBadRequest)
			logError(err)
			return
		}
		record, getErr = a.records.GetRecordVersion(ctx, int(idNumber), int(versionNumber))
	}

	if getErr != nil {
		err := writeError(w, fmt.Sprintf("failed to get record: %v", getErr), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}

func (a *API) GetRecordVersionsV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)
	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	versions, err := a.records.GetRecordVersions(ctx, int(idNumber))
	if err != nil {
		err := writeError(w, fmt.Sprintf("failed to get record versions: %v", err), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, versions, http.StatusOK)
	logError(err)
}
