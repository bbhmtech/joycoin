package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bbhmtech/joycoin/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"
)

type ctxKey int

const _ctxSessionAccount ctxKey = 0

type APIServerV1 struct {
	db  *gorm.DB
	scc *securecookie.SecureCookie
}

func (s *APIServerV1) SessionAccount(r *http.Request) *model.Account {
	return r.Context().Value(_ctxSessionAccount).(*model.Account)
}

func (s *APIServerV1) writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (s *APIServerV1) readJSON(r *http.Request) (map[string]interface{}, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *APIServerV1) AccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	acc := model.Account{ID: uint(accId)}
	err = s.db.Take(&acc).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, acc)
}

func (s *APIServerV1) TransactionActionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessAcc := s.SessionAccount(r)
	switch r.Method {
	case http.MethodGet:
		tID, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "incorrect transaction ID", http.StatusBadRequest)
			return
		}

		t := model.Transaction{ID: uint(tID)}
		err = s.db.Take(&t).Error
		if err != nil {
			http.Error(w, "cannot lookup tID in database", http.StatusNotFound)
			return
		}

		s.writeJSON(w, t)
		return
	case http.MethodPut:
		// ID 随便填吧
		data, err := s.readJSON(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t := model.Transaction{
			ReferenceTag:       data["reference_tag"].(string),
			InitiatorAccountID: sessAcc.ID,
			FromAccountID:      uint(data["from"].(float64)),
			ToAccountID:        uint(data["from"].(float64)),
		}

		if err = t.PreFlightCheck(s.db); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = s.db.Create(&t).Error; err != nil {
			http.Error(w, "交易被数据库拒绝: "+err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
			// TODO: 返回真实 ID
		}
		return
	case http.MethodDelete:
		http.Error(w, "TODO: 撤回交易", http.StatusNotImplemented)
		return
	}
}

func (s *APIServerV1) ListTransactionHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *APIServerV1) authRequired(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sccv := decodeSecureCookie(w, r, s.scc)
		sessionAcc, err := sccv.GetAccount(s.db)
		loggedIn := sessionAcc != nil
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !loggedIn {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), _ctxSessionAccount, sessionAcc)

		handler(w, r.WithContext(ctx))
	})
}
func CreateAPIServerV1(db *gorm.DB, scc *securecookie.SecureCookie) {
	s := APIServerV1{db, scc}
	r := http.ServeMux{}
	r.Handle("/_/v1/account/{id}", s.authRequired(s.AccountHandler))

	r.Handle("/_/v1/transaction/{id}", s.authRequired(s.TransactionActionHandler))
	r.Handle("/_/v1/transaction", s.authRequired(s.ListTransactionHandler))
}
