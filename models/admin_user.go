package models

import (
	"database/sql"

	"../config"
)

//GetPasswordFromEmail emailからパスワードを取得する
func GetPasswordFromEmail(email string) (string, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	var password string
	err = db.QueryRow("SELECT password FROM admin_user WHERE email = ?", email).Scan(&password)
	return password, err
}
