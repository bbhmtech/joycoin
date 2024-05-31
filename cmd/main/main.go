package main

import (
	"net/http"

	"github.com/bbhmtech/joycoin"
	"github.com/bbhmtech/joycoin/api"
	"github.com/bbhmtech/joycoin/model"
)

func main() {
	cfg := joycoin.LoadConfig("config.json")
	db := cfg.InitializeDatabase()
	secc := cfg.InitializeSecureCookie()
	model.AutoMigration(db)

	http.Handle("/j/", api.CreateJumperServer(db, secc, cfg))
	http.Handle("/_/v1/", api.CreateAPIServerV1(db, secc, cfg))
	fs := http.FileServer(http.FS(joycoin.StaticContent))
	http.Handle("/", fs)

	// KeyPairWithPin()
	// http.ListenAndServeTLS(cfg.ListenAddr, "cert.pem", "key.pem", nil)

	http.ListenAndServe(cfg.ListenAddr, nil)
	// fmt.Println(model.JumperMapFromKey("S26awVL98GSvZp15wsJQ9Q"))
}
