package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/guoxiaopeng875/wallet/internal/pkg/util"
	"net/http"
)

// Helper functions for request handling

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			handleError(w, errors.InternalServer.WithCause(err))
		}
	}
}

func parseReqBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		handleError(w, errors.InvalidArgs.WithCause(err))
		return false
	}
	return true
}

func parseWalletID(w http.ResponseWriter, r *http.Request) uint {
	id, err := util.StringToUint(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, errors.InvalidArgs.WithCause(err))
		return 0
	}
	return id
}

func handleError(w http.ResponseWriter, err error) {
	var wErr *errors.Error
	if errors.As(err, &wErr) {
		http.Error(w, wErr.Message, wErr.Code)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
