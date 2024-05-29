package api

import (
	"net/http"
	"time"

	"github.com/bbhmtech/joycoin"
	"github.com/bbhmtech/joycoin/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"
)

type JumperServer struct {
	db         *gorm.DB
	scc        *securecookie.SecureCookie
	qpRedirect string
}

func (s *JumperServer) handleAccount(w http.ResponseWriter, r *http.Request, j *model.Jumper) {
	sccv := decodeSecureCookie(w, r, s.scc)
	qa, err := sccv.GetQuickAction(s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tagAcc := model.Account{ID: j.TargetID}
	err = s.db.Take(&tagAcc).Error
	if err == gorm.ErrRecordNotFound {
		// bad tag
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !tagAcc.Activated {
		// not activated account, redirect to activate confirmation.
		// additional hint if already loggedIn
		return
	}

	if qa != nil {
		switch qa.Action {
		case "quickPay":
			err = s.db.Transaction(func(tx *gorm.DB) error {
				var err error
				var t model.Transaction
				if qa.Int64Value1 > 0 {
					t = model.Transaction{
						ReferenceTag:       "quickPay" + time.Now().Format(time.RFC3339),
						InitiatorAccountID: qa.CachedAccountID,
						FromAccountID:      qa.CachedAccountID,
						ToAccountID:        tagAcc.ID,
						CentAmount:         qa.Int64Value1,
						Message:            qa.StringValue1,
					}
				} else {
					t = model.Transaction{
						ReferenceTag:       "quickPay" + time.Now().Format(time.RFC3339),
						InitiatorAccountID: qa.CachedAccountID,
						FromAccountID:      tagAcc.ID,
						ToAccountID:        qa.CachedAccountID,
						CentAmount:         qa.Int64Value1,
						Message:            qa.StringValue1,
					}
				}

				if err = t.PreFlightCheck(tx); err != nil {
					return err
				}

				if err = t.Insert(tx); err != nil {
					return err
				}

				if qa.Temporary {
					qa.ValidBefore = time.Now()
					err = tx.Save(&qa).Error
				}
				return err
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// TODO: pass transaction details
			http.Redirect(w, r, s.qpRedirect, http.StatusTemporaryRedirect)
			return
		default:
			http.Error(w, "unsupported action", http.StatusInternalServerError)
			return
		}
	}

	sessAcc, err := sccv.GetAccount(s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sessAcc != nil {
		// http.Redirect() to account page
		return
	}

	// redirect to login page
}

func (s *JumperServer) handleSLink(w http.ResponseWriter, r *http.Request, j *model.Jumper) {
	slink := model.ShortenLink{ID: j.TargetID}
	s.db.Take(&slink)
	http.Redirect(w, r, slink.TargetURL, http.StatusTemporaryRedirect)
	// http.Redirect(w, r, j.TargetURL, http.StatusPermanentRedirect)
}

func (s *JumperServer) HandleJ(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	j, err := model.JumperFromEncodedID(s.db, vars["key"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch j.Hint {
	case "NTAG|Account":
		s.handleAccount(w, r, j)
	case "NTAG|SLink":
		s.handleSLink(w, r, j)
	default:
		http.Error(w, "unknown j.Hint="+j.Hint, http.StatusInternalServerError)
		return
	}
}

func CreateJumperServer(db *gorm.DB, secc *securecookie.SecureCookie, cfg *joycoin.Config) http.Handler {
	j := JumperServer{
		db:         db,
		scc:        secc,
		qpRedirect: cfg.QuickPayResultURL,
	}

	r := mux.NewRouter()
	r.HandleFunc("/j/{key}", j.HandleJ)
	return r
}
