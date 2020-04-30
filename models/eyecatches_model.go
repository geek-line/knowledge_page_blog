package models

import (
	"database/sql"

	"../config"
	"../structs"
)

//GetAllEyecatches eyecatchを全て取得する
func GetAllEyecatches() ([]structs.Eyecatch, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("SELECT id, name, src FROM eyecatches")
	defer rows.Close()
	var eyecatches []structs.Eyecatch
	for rows.Next() {
		var eyecatch structs.Eyecatch
		err = rows.Scan(&eyecatch.ID, &eyecatch.Name, &eyecatch.Src)
		eyecatches = append(eyecatches, eyecatch)
	}
	return eyecatches, err
}

//PostEyecatch eyecatchをsbに追加する
func PostEyecatch(name string, src string) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("INSERT INTO eyecatches(name, src) VALUES(?, ?)", name, src)
	defer rows.Close()
	return err
}

//UpdateEyecatch 指定されたidのeyecatchを更新する
func UpdateEyecatch(id int, name string, src string) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("UPDATE eyecatches SET name = ?, src = ? WHERE id = ?", name, src, id)
	defer rows.Close()
	return err
}

//DeleteEyecatch 指定されたidのeyecatchを削除する
func DeleteEyecatch(id int) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("DELETE FROM eyecatches WHERE id = ?", id)
	defer rows.Close()
	return err
}
