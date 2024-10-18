package main

import (
	"github.com/mytempoesp/mysql-easy"
)

func (a *Ay) SetupDatabase() (err error) {

	db, err := mysql_easy.ConfiguraDB()

	a.db = db

	return
}

func (a *Ay) CloseDatabase() {

	a.db.Close()
}
