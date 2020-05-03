package models

import (
	"database/sql"
	"time"

	"../config"
	"../structs"
)

//GetAllTags 全てのtagを取得する
func GetAllTags() ([]structs.Tag, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("SELECT id, name FROM tags")
	defer rows.Close()
	var tags []structs.Tag
	for rows.Next() {
		var tag structs.Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		tags = append(tags, tag)
	}
	return tags, err
}

//GetTagFromKnowledgeID 指定されたknowledgeのidからついているタグを取得する
func GetTagFromKnowledgeID(id int) ([]structs.Tag, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	tagsRows, err := db.Query("SELECT tags.id, tags.name FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id WHERE knowledge_id = ?", id)
	defer tagsRows.Close()
	var selectedTags []structs.Tag
	for tagsRows.Next() {
		var selectedTag structs.Tag
		err = tagsRows.Scan(&selectedTag.ID, &selectedTag.Name)
		selectedTags = append(selectedTags, selectedTag)
	}
	return selectedTags, err
}

//PostTag tagを新規作成する
func PostTag(name string, createdAt time.Time, updatedAt time.Time) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("INSERT INTO tags(name, created_at, updated_at) VALUES(?, ?, ?)", name, createdAt, updatedAt)
	defer rows.Close()
	return err
}

//UpdateTag 指定されたidのtagを更新する
func UpdateTag(id int, name string, updatedAt time.Time) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("UPDATE tags SET name = ?, updated_at = ? WHERE id = ?", name, updatedAt, id)
	defer rows.Close()
	return err
}

//DeleteTag 指定されたidのtagを削除する
func DeleteTag(id int) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("DELETE FROM tags WHERE id = ?", id)
	defer rows.Close()
	return err
}

//GetTop10ReferencedTags 被参照数が多い上位10のtagを取得する
func GetTop10ReferencedTags() ([]structs.TagRankingElem, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("SELECT tags.id, tags.name, count(knowledges_tags.id) AS count FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id GROUP BY knowledges_tags.tag_id ORDER BY count DESC LIMIT 10;")
	defer rows.Close()
	var tags []structs.TagRankingElem
	for rows.Next() {
		var tag structs.TagRankingElem
		err = rows.Scan(&tag.TagID, &tag.TagName, &tag.CountOfRefferenced)
		tags = append(tags, tag)
	}
	return tags, err
}
