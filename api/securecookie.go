package api

import (
	"net/http"

	"github.com/bbhmtech/joycoin/model"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"
)

const _secureCookieName = "secure-joycoin"

type mySecureCookieValue struct {
	rawValue map[string]string
}

func decodeSecureCookie(w http.ResponseWriter, r *http.Request, scc *securecookie.SecureCookie) *mySecureCookieValue {
	if cookie, err := r.Cookie(_secureCookieName); err == nil {
		rawValue := make(map[string]string)
		if err = scc.Decode(_secureCookieName, cookie.Value, &rawValue); err == nil {
			return &mySecureCookieValue{rawValue}
		} else {
			cookie := &http.Cookie{
				Name:   _secureCookieName,
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)
		}
	}
	return nil
}

func (v *mySecureCookieValue) SetCookie(w http.ResponseWriter, scc *securecookie.SecureCookie) {
	if encoded, err := scc.Encode(_secureCookieName, v.rawValue); err == nil {
		cookie := &http.Cookie{
			Name:     _secureCookieName,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   86400,
		}
		http.SetCookie(w, cookie)
	}
}

func (v *mySecureCookieValue) DeviceBindingKey() string {
	return v.rawValue["dbk"]
}

func (v *mySecureCookieValue) GetAccount(db *gorm.DB) (*model.Account, error) {
	clientKey := v.DeviceBindingKey()
	if len(clientKey) > 0 {
		acc := model.Account{DeviceBindingKey: clientKey}
		err := db.Take(&acc).Error
		if err == nil {
			return &acc, nil
		} else if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}

	}
	return nil, nil
}
