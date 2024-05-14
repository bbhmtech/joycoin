package api

import (
	"net/http"

	"github.com/bbhmtech/joycoin/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"
)

type JumperServer struct {
	db  *gorm.DB
	scc *securecookie.SecureCookie
}

func (s *JumperServer) handleAccount(w http.ResponseWriter, r *http.Request, j *model.Jumper) {
	sccv := decodeSecureCookie(w, r, s.scc)
	sessionAcc, err := sccv.GetAccount(s.db)
	loggedIn := sessionAcc != nil
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

	if loggedIn {
		if sessionAcc.Role == "merchant" {
			if tagAcc.Role == "normal" {
				// do quickAction
			} else if tagAcc.ID == sessionAcc.ID {
				// merchant dashboard
			}
		} else if sessionAcc.Role == "normal" && tagAcc.Role == "normal" {
			if sessionAcc.ID == tagAcc.ID {
				// normal dashboard
			} else {
				// check oneshot pay
			}
		}
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

func CreateJumperServer(db *gorm.DB, secc *securecookie.SecureCookie) http.Handler {
	j := JumperServer{db, secc}

	r := mux.NewRouter()
	r.HandleFunc("/j/{key}", j.HandleJ)
	return r
}
