package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bbhmtech/joycoin/model"
	"github.com/google/uuid"
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
	accID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessAcc := s.SessionAccount(r)
	//
	switch r.Method {
	case http.MethodGet:
		if !sessAcc.IsOperator() && sessAcc.ID != uint(accID) {
			http.Error(w, "未授权", http.StatusForbidden)
			return
		}
		acc := model.Account{ID: uint(accID)}
		err = s.db.Take(&acc).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		s.writeJSON(w, acc)
	case http.MethodPost:
		if !sessAcc.IsOperator() {
			http.Error(w, `"现在还用不了！"`, http.StatusForbidden)
			return
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		acc := model.Account{}
		err = json.Unmarshal(b, &acc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.db.Save(&acc).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServerV1) AccountActivateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	data, err := s.readJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if acc.Activated {
		// already activated, match passcode
		if !acc.VerifyPasscode(data["passcode"].(string)) {
			http.Error(w, "wrong password", http.StatusForbidden)
			return
		}
	} else {
		// not activated, set passcode
		acc.ChangePasscode(data["passcode"].(string))
		acc.Activated = true
		if err := s.db.Model(&acc).Where("activated = ?", false).Select("PasscodeHash", "Activated").Updates(&acc).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// login if no problem
	u := uuid.New()
	acc.DeviceBindingKey = u.String()
	if err := s.db.Model(&acc).Where("activated = ?", true).Select("DeviceBindingKey").Updates(&acc).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	createSecureCookieValue(u).SetCookie(w, s.scc)
	w.WriteHeader(http.StatusOK)
}

// princple here: 请求必须等幂
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
			ToAccountID:        uint(data["to"].(float64)),
			CentAmount:         int64(data["cent_amount"].(float64)),
			Message:            data["message"].(string),
		}

		if err = t.PreFlightCheck(s.db); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = t.Insert(s.db); err != nil {
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
func CreateAPIServerV1(db *gorm.DB, scc *securecookie.SecureCookie) http.Handler {
	s := APIServerV1{db, scc}
	r := mux.NewRouter()
	r.HandleFunc("/_/v1/account/{id}/activate", s.AccountActivateHandler)
	r.Handle("/_/v1/account/{id}", s.authRequired(s.AccountHandler))

	r.Handle("/_/v1/transaction/{id}", s.authRequired(s.TransactionActionHandler))
	r.Handle("/_/v1/transaction", s.authRequired(s.ListTransactionHandler))
	return r
}
