package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/bbhmtech/joycoin"
	"github.com/bbhmtech/joycoin/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"
)

type APIServerV1 struct {
	db         *gorm.DB
	scc        *securecookie.SecureCookie
	corsOrigin string
	jURLPrefix string
}

var _success = map[string]bool{"ok": true}

func (s *APIServerV1) writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	accID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessAcc := sessionAccount(r)
	if accID != 0 {
		if !sessAcc.IsOperator() && sessAcc.ID != uint(accID) {
			http.Error(w, "未授权", http.StatusForbidden)
			return
		}

	} else {
		accID = uint64(sessAcc.ID)
	}

	switch r.Method {
	case http.MethodGet:
		acc := model.Account{ID: uint(accID)}
		err = s.db.Take(&acc).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		s.writeJSON(w, acc)
	case http.MethodPost:
		data, err := s.readJSON(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newAcc := model.Account{ID: uint(accID), Nickname: data["nickname"].(string)}
		switch data["passcode"].(type) {
		case string:
			pass := data["passcode"].(string)
			if len(pass) > 0 {
				newAcc.ChangePasscode(pass)
			}
		}

		err = s.db.Model(&sessAcc).Select("Nickname", "PasscodeHash").Updates(newAcc).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, _success)
	}
}

func (s *APIServerV1) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	sessAcc := sessionAccount(r)
	if !sessAcc.IsOperator() {
		http.Error(w, "operator only", http.StatusForbidden)
		return
	}

	data, err := s.readJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a := model.Account{Role: data["role"].(string), Activated: false}
	if data["create_jumper"].(bool) {
		j := model.Jumper{ID: uuid.NewString(), Hint: "NTAG|Account"}
		err = s.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Create(&a).Error
			if err != nil {
				return err
			}
			j.TargetID = a.ID

			err = tx.Create(&j).Error
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		eID, err := j.EncodeID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, map[string]interface{}{
			"account": a,
			"jumper":  j,
			"link":    s.jURLPrefix + eID,
		})

	} else {
		err := s.db.Create(&a).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, a)
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

		acc.Nickname = data["nickname"].(string)
		acc.Activated = true
		if err := s.db.Model(&acc).Where("activated = ?", false).Select("PasscodeHash", "Nickname", "Activated").Updates(&acc).Error; err != nil {
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

	s.writeJSON(w, _success)
}

// princple here: 请求必须等幂
func (s *APIServerV1) TransactionActionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessAcc := sessionAccount(r)
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
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}

	sessAcc := sessionAccount(r)
	ts := []model.Transaction{}
	err := s.db.Where("from_account_id = ? or to_account_id = ?", sessAcc.ID, sessAcc.ID).Order("updated_at DESC").Find(&ts).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, ts)
}

func (s *APIServerV1) QuickActionHandler(w http.ResponseWriter, r *http.Request) {
	sessAcc := sessionAccount(r)
	switch r.Method {
	case http.MethodPost:
		{
			data, err := s.readJSON(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			switch data["action"].(string) {
			case "quickPay":
				amount := int64(data["cent_amount"].(float64))
				// amount > 0: pay from sessAcc to ...
				// amount < 0: collect money from ... to sessAcc
				if amount < 0 && sessAcc.IsNormal() {
					http.Error(w, "普通账户不支持快速收款", http.StatusForbidden)
					return
				}

				if data["temporary"].(bool) {
					qa := model.QuickAction{
						DeviceBindingKey: sessAcc.DeviceBindingKey,
						ValidBefore:      time.Now().Add(time.Minute),
						Temporary:        true,
						CachedAccountID:  sessAcc.ID,
						Action:           "quickPay",
						Int64Value1:      amount,
						StringValue1:     data["message"].(string),
					}
					s.db.Save(&qa)
				} else {
					qa := model.QuickAction{
						DeviceBindingKey: sessAcc.DeviceBindingKey,
						ValidBefore:      time.Now().Add(24 * time.Hour),
						Temporary:        false,
						CachedAccountID:  sessAcc.ID,
						Action:           "quickPay",
						Int64Value1:      amount,
						StringValue1:     data["message"].(string),
					}
					s.db.Save(&qa)
				}
			case "null":
				fallthrough
			default:
				qa := model.QuickAction{
					DeviceBindingKey: sessAcc.DeviceBindingKey,
					ValidBefore:      time.Now(),
					Temporary:        false,
					CachedAccountID:  sessAcc.ID,
					Action:           "null",
				}

				ret := s.db.Model(&qa).Select("valid_before", "temporary", "action").Updates(&qa)
				slog.Debug("saved null quickAction", "err", ret.Error, "rowAffected", ret.RowsAffected)
			}

			s.writeJSON(w, _success)
			return
		}
	case http.MethodGet:
		{
			qa := model.QuickAction{
				DeviceBindingKey: sessAcc.DeviceBindingKey,
			}

			ret := s.db.Joins("CachedAccount").Limit(1).Find(&qa)
			if ret.Error != nil {
				http.Error(w, ret.Error.Error(), http.StatusInternalServerError)
				return
			}
			if ret.RowsAffected == 0 || !qa.IsValid() {
				http.Error(w, "no valid QuickAction found", http.StatusNotFound)
				return
			}

			s.writeJSON(w, qa)
		}
	default:
		{
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func (s *APIServerV1) CORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (s *APIServerV1) ListJumpers(w http.ResponseWriter, r *http.Request) {
	sessAcc := sessionAccount(r)
	if sessAcc.IsOperator() {
		var data []model.Jumper
		err := s.db.Find(&data).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, data)
	} else {
		http.Error(w, "OP only", http.StatusForbidden)
	}
}

func CreateAPIServerV1(db *gorm.DB, scc *securecookie.SecureCookie, cfg *joycoin.Config) http.Handler {
	s := APIServerV1{
		db:         db,
		scc:        scc,
		corsOrigin: cfg.AllowedCORSOrigin,
		jURLPrefix: cfg.JumperURLPrefix,
	}

	r := mux.NewRouter()
	r.HandleFunc("/_/v1/account/{id:[0-9]+}/activate", s.AccountActivateHandler)
	r.Handle("/_/v1/account/{id:[0-9]+}", authRequired(s.scc, s.db, s.AccountHandler))
	r.Handle("/_/v1/account/create", authRequired(s.scc, s.db, s.CreateAccountHandler))
	r.Handle("/_/v1/quickaction", authRequired(s.scc, s.db, s.QuickActionHandler))

	r.Handle("/_/v1/transaction/{id:[0-9]+}", authRequired(s.scc, s.db, s.TransactionActionHandler))
	r.Handle("/_/v1/transaction", authRequired(s.scc, s.db, s.ListTransactionHandler))

	r.Handle("/_/v1/jumper", authRequired(s.scc, s.db, s.ListJumpers))

	r.Use(s.CORSHeaders)
	return r
}
