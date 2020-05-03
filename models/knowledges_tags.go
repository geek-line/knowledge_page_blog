package models

import (
	"database/sql"
	"time"

	"../config"
)

//DeleteKnowledgesTagsFromKnowledgeID 指定されたknowledgeのidのknowledges_tagsを削除する
func DeleteKnowledgesTagsFromKnowledgeID(id int) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("DELETE FROM knowledges_tags WHERE knowledge_id = ?", id)
	defer rows.Close()
	return err
}

//GetTagIDsFromKnowledgeID 指定されたknowledgeのidからtagのidを取得する
func GetTagIDsFromKnowledgeID(id int) ([]int, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	var selectedTagIDs []int
	rows, _ := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", id)
	for rows.Next() {
		var selectedTagID int
		err = rows.Scan(&selectedTagID)
		selectedTagIDs = append(selectedTagIDs, selectedTagID)
	}
	return selectedTagIDs, err
}

//PostKnowledgesTags KnowledgesTagsを新規作成する
func PostKnowledgesTags(knowledgeID int, tagID int, createdAt time.Time, updatedAt time.Time) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("INSERT INTO knowledges_tags(knowledge_id, tag_id, created_at, updated_at) VALUES(?, ?, ?, ?)", knowledgeID, tagID, createdAt, updatedAt)
	defer rows.Close()
	return err
}

//DeleteKnowledgesTagsFromKnowledgeIDFromTagID 指定されたtagのidからknowledges_tagsを削除する
func DeleteKnowledgesTagsFromKnowledgeIDFromTagID(id int) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("DELETE FROM knowledges_tags WHERE tag_id = ?", id)
	defer rows.Close()
	return err
}
