package api

import (
	"net/http"
)

type APIServerV1 struct {
}

func (s *APIServerV1) AccountHandler(w http.ResponseWriter, r *http.Request) {

}

func CreateAPIServerV1() {
	s := APIServerV1{}
	r := http.ServeMux{}
	r.HandleFunc("/_/v1/account", s.AccountHandler)
}
