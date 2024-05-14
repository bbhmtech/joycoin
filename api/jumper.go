package api

import (
	"net/http"

	"github.com/bbhmtech/joycoin/model"
	"github.com/gorilla/mux"
)

type JumperServer struct {
}

func (s *JumperServer) HandleJ(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	j, err := model.JumperMapFromEncodedID(vars["key"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, j.TargetURL, http.StatusTemporaryRedirect)
	// http.Redirect(w, r, j.TargetURL, http.StatusPermanentRedirect)
	return
}

func CreateJumperServer() http.Handler {
	j := JumperServer{}

	r := mux.NewRouter()
	r.HandleFunc("/j/{key}", j.HandleJ)
	return r
}
