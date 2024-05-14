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

	http.Handle("/j/", api.CreateJumperServer(db, secc))
	http.ListenAndServe(":8080", nil)
	// fmt.Println(model.JumperMapFromKey("S26awVL98GSvZp15wsJQ9Q"))
}
