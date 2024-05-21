package api

import (
	"net/http"

	"github.com/bbhmtech/joycoin/model"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"
)

const _secureCookieName = "secure-joycoin-v1"
const _sccDeviceBindingKey = "dbk"

type mySecureCookieValue struct {
	rawValue map[string]string
}

func createSecureCookieValue(deviceBindingKey string) *mySecureCookieValue {
	v := map[string]string{
		_sccDeviceBindingKey: deviceBindingKey,
	}
	return &mySecureCookieValue{
		rawValue: v,
	}
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
	encoded, err := scc.Encode(_secureCookieName, v.rawValue)
	if err == nil {
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

func (v *mySecureCookieValue) GetAccount(db *gorm.DB) (*model.Account, error) {
	if v == nil {
		return nil, nil
	}

	clientKey := v.rawValue[_sccDeviceBindingKey]
	if len(clientKey) > 0 {
		acc := model.Account{DeviceBindingKey: clientKey}
		err := db.Where(&acc).First(&acc).Error
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
