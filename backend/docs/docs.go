package docs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"task-manager-api/db"
	"task-manager-api/task"
)

type Document struct {
	ID        int    `json:"Docs_ID"`
	// task_ID   int    `json:"Task_ID"`
	Doc       string `json:"content"`
  Create_At string `json:"Create_At"`
}

// var documents = make(map[int][]Document)

func GetDocuments(w http.ResponseWriter, r *http.Request){
	_, Str, err := task.ExtractUserIDAndParam(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}
	fmt.Println("idStr : ",Str)

	// URLパスからパラメータを取得
	strid := strings.TrimPrefix(Str, "documents/")
	
	fmt.Println("idpara : ",strid)

	taskid, err := strconv.Atoi(strid)
	fmt.Println("id : ",taskid)
	if err != nil {
		http.Error(w, "無効なIDです1", http.StatusBadRequest)
		return
	}
  
  rows, err := db.DB.Query("SELECT Documents FROM Docs WHERE Task_ID= ?",taskid)
	if err != nil {
		http.Error(w, `{"error":"データベースエラー"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.Doc)
		if err != nil {
			http.Error(w,`{"error": "データ取得エラー"}`, http.StatusInternalServerError)
			return
		}
		docs = append(docs, doc)
	}

  fmt.Println("docs: ",docs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func AddDocument(w http.ResponseWriter, r *http.Request) {
	_, Str, err := task.ExtractUserIDAndParam(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}
	fmt.Println("idStr : ",Str)

	// URLパスからパラメータを取得
	strid := strings.TrimPrefix(Str, "documents/")
	
	fmt.Println("idpara : ",strid)

	taskID, err := strconv.Atoi(strid)
	fmt.Println("id : ",taskID)
	if err != nil {
		http.Error(w, "無効なIDです1", http.StatusBadRequest)
		return
	}

  var NewDoc Document
	err = json.NewDecoder(r.Body).Decode(&NewDoc)
	if err != nil {
    fmt.Print("JSONデコードエラー")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentTime,err := task.MakeTimeStamp()
	if err != nil {
		fmt.Println("Error generating timestamp:", err)
		return
	}

  fmt.Println("Documents: ",NewDoc.Doc)
  fmt.Println("Documents: ",NewDoc)

	// データベースにタスクを追加
	_, err = db.DB.Exec("INSERT INTO Docs (Task_ID, Documents, Create_At) VALUES (?, ?, ?)", taskID, NewDoc.Doc, currentTime)
	if err != nil {
		http.Error(w, `{"error":"データベースエラー"}`, http.StatusInternalServerError)
		return
	}






	
}