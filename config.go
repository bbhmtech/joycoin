package joycoin

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/gorilla/securecookie"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	_fromFileName           string `json:"-"`
	DatabseConnectionString string
	SecureCookieHashKey     string
	SecureCookieBlockKey    string
	ListenAddr              string
	AllowedCORSOrigin       string
}

func LoadConfig(filename string) *Config {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	cfg := Config{_fromFileName: filename}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func (c *Config) SaveConfig() {
	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(c._fromFileName, b, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (c *Config) InitializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(c.DatabseConnectionString))
	if err != nil {
		panic(err)
	}

	return db
}

func (c *Config) InitializeSecureCookie() *securecookie.SecureCookie {
	if len(c.SecureCookieBlockKey) == 0 || len(c.SecureCookieHashKey) == 0 {
		c.SecureCookieHashKey = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64))
		c.SecureCookieBlockKey = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
		c.SaveConfig()
	}
	hKey, err := base64.StdEncoding.DecodeString(c.SecureCookieHashKey)
	if err != nil {
		panic(err)
	}
	bKey, err := base64.StdEncoding.DecodeString(c.SecureCookieBlockKey)
	if err != nil {
		panic(err)
	}

	s := securecookie.New(hKey, bKey)

	return s
}
