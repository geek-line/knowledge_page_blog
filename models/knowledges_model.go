package models

import (
	"database/sql"
	"fmt"
	"time"

	"../config"
	"../structs"
)

//DeleteKnowledge 指定されたidのknowledgeを削除する
func DeleteKnowledge(id int) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("DELETE FROM knowledges WHERE id = ?", id)
	defer rows.Close()
	return err
}

//GetKnowledge 指定されたidのknowledgeを取得する
func GetKnowledge(id int) (structs.Knowledge, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	var knowledge structs.Knowledge
	err = db.QueryRow("SELECT id, title, content, updated_at, likes, eyecatch_src FROM knowledges WHERE id = ?", id).Scan(&knowledge.ID, &knowledge.Title, &knowledge.Content, &knowledge.UpdatedAt, &knowledge.Likes, &knowledge.EyecatchSrc)
	return knowledge, err
}

//GetAllKnowledges 全てのknowledgeを取得する
func GetAllKnowledges() ([]structs.Knowledge, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	var knowledges []structs.Knowledge
	rows, err := db.Query("SELECT id, title, created_at, updated_at FROM knowledges")
	defer rows.Close()
	for rows.Next() {
		var knowledge structs.Knowledge
		err = rows.Scan(&knowledge.ID, &knowledge.Title, &knowledge.CreatedAt, &knowledge.UpdatedAt)
		knowledges = append(knowledges, knowledge)
	}
	return knowledges, err
}

//PostKnowledge knowledgeを新規作成して作成したknowledgeのIDを取得する
func PostKnowledge(title string, content string, rowContent string, createdAt time.Time, updatedAt time.Time, eyecatchSrc string) (int64, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	stmtInsert, err := db.Prepare("INSERT INTO knowledges(title, content, row_content, created_at, updated_at, eyecatch_src) VALUES(?, ?, ?, ?, ?, ?)")
	defer stmtInsert.Close()
	result, err := stmtInsert.Exec(title, content, rowContent, createdAt, updatedAt, eyecatchSrc)
	knowledgeID, err := result.LastInsertId()
	return knowledgeID, err
}

//UpdateKnowledge 指定したidのknowledgeを更新する
func UpdateKnowledge(knowledgeID int, title string, content string, rowContent string, updatedAt time.Time, eyecatchSrc string) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("UPDATE knowledges SET title = ?, content = ?, row_content = ?, updated_at = ?, eyecatch_src = ? WHERE id = ?", title, content, rowContent, updatedAt, eyecatchSrc, knowledgeID)
	defer rows.Close()
	return err
}

//IncrementLikes 指定されたidのlikesを増やす
func IncrementLikes(id int) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("UPDATE knowledges SET likes = likes + 1 WHERE id = ?", id)
	defer rows.Close()
	return err
}

//DecrementLikes 指定されたidのlikesを減らす
func DecrementLikes(id int) error {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	rows, err := db.Query("UPDATE knowledges SET likes = likes - 1 WHERE id = ?", id)
	defer rows.Close()
	return err
}

//GetNumOfKnowledges knowledgeの数を取得する
func GetNumOfKnowledges() (float64, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	var numOfKnowledges float64
	db.QueryRow("SELECT count(id) FROM knowledges").Scan(&numOfKnowledges)
	return numOfKnowledges, err
}

//GetNumOfKnowledgesFilteredTagID 指定されたtag_idに該当するknowledgeの数を取得する
func GetNumOfKnowledgesFilteredTagID(id int) (float64, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	var numOfKnowledges float64
	db.QueryRow("SELECT count(knowledge_id) FROM knowledges_tags WHERE tag_id = ?", id).Scan(&numOfKnowledges)
	return numOfKnowledges, err
}

//Get20SortedElems 指定のsortKeyでソートされた20のknowledgeの要素を取得する
func Get20SortedElems(sortKey string, startIndex int, length int) ([]structs.IndexElem, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	qtext := fmt.Sprintf("SELECT id, title, updated_at, likes, eyecatch_src FROM knowledges ORDER BY %s DESC LIMIT ?, ?", sortKey)
	rows, err := db.Query(qtext, startIndex, length)
	defer rows.Close()
	var indexElems []structs.IndexElem
	for rows.Next() {
		var indexElem structs.IndexElem
		err = rows.Scan(&indexElem.Knowledge.ID, &indexElem.Knowledge.Title, &indexElem.Knowledge.UpdatedAt, &indexElem.Knowledge.Likes, &indexElem.Knowledge.EyecatchSrc)
		var selectedTags []structs.Tag
		var tagsRows *sql.Rows
		tagsRows, err = db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", indexElem.Knowledge.ID)
		defer tagsRows.Close()
		for tagsRows.Next() {
			var selectedTag structs.Tag
			err = tagsRows.Scan(&selectedTag.ID)
			db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTag.ID).Scan(&selectedTag.Name)
			selectedTags = append(selectedTags, selectedTag)
		}
		indexElem.SelectedTags = selectedTags
		indexElems = append(indexElems, indexElem)
	}
	return indexElems, err
}

//Get20SortedElemFilteredTagID 指定のsortKeyでソートされ、指定のtagIdでフィルターされた20のknowledgeの要素を取得する
func Get20SortedElemFilteredTagID(sortKey string, tagID int, startIndex int, length int) ([]structs.IndexElem, string, error) {
	db, err := sql.Open("mysql", config.SQLEnv)
	defer db.Close()
	var tagName string
	if err = db.QueryRow("SELECT name FROM tags WHERE id = ?", tagID).Scan(&tagName); err != nil {
		return nil, "", err
	}
	qtext := fmt.Sprintf("SELECT knowledges.id, title, knowledges.updated_at, likes, eyecatch_src FROM knowledges INNER JOIN knowledges_tags ON knowledges_tags.knowledge_id = knowledges.id WHERE tag_id = ? ORDER BY %s DESC LIMIT ?, ?", "updated_at")
	rows, err := db.Query(qtext, tagID, startIndex, length)
	defer rows.Close()
	var indexElems []structs.IndexElem
	for rows.Next() {
		var indexElem structs.IndexElem
		err = rows.Scan(&indexElem.Knowledge.ID, &indexElem.Knowledge.Title, &indexElem.Knowledge.UpdatedAt, &indexElem.Knowledge.Likes, &indexElem.Knowledge.EyecatchSrc)
		var selectedTags []structs.Tag
		var tagsRows *sql.Rows
		tagsRows, err = db.Query("SELECT tags.id, tags.name FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id WHERE knowledge_id = ?", indexElem.Knowledge.ID)
		defer tagsRows.Close()
		for tagsRows.Next() {
			var selectedTag structs.Tag
			err = tagsRows.Scan(&selectedTag.ID, &selectedTag.Name)
			selectedTags = append(selectedTags, selectedTag)
		}
		indexElem.SelectedTags = selectedTags
		indexElems = append(indexElems, indexElem)
	}
	return indexElems, tagName, err
}
