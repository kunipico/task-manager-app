package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"task-manager-api/db"

	"github.com/golang-jwt/jwt"
)

// タスク構造体
type Task struct {
	ID   int    `json:"Task_ID"`
	Name string `json:"Task_Name"`
	Done string   `json:"Task_Done"`
}

// JWTシークレットキー
var jwtKey = []byte("my_secret_key")

// JWTのペイロードを定義
type Claims struct {
	Userinfo string `json:"userinfo"`
	jwt.StandardClaims
}

//Task一覧取得処理
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Cookieの内容を確認
	cookie, err := r.Cookie("token")
	fmt.Println("Cookie : ", cookie)
  if err != nil {
    http.Error(w, `{"error":"Nothing Cookie!"}`, http.StatusUnauthorized)
		// http.Error(w,"Unauthorized", http.StatusUnauthorized)
		fmt.Println("Nothing Cookie!!!")
    return
  }
  tokenStr := cookie.Value
  // JWTの検証
  claims := &Claims{}
  token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
    return jwtKey, nil
  })
  if err != nil || !token.Valid {
    http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
    return
  }

  // クレームからユーザーIDを取得
  userID := claims.Userinfo
	fmt.Println("userID : ",userID)
  // DBからユーザー情報を取得
  // var username, email string
  // err = db.QueryRow("SELECT username, email FROM Users WHERE user_id = ?", userID).Scan(&username, &email)
  // if err != nil {
  //   http.Error(w, "User not found", http.StatusNotFound)
  //   return
  // }

  // // ユーザー情報をJSON形式で返す
  // response := map[string]string{
  //   "username": username,
  //   "email":    email,
  // }
  // w.Header().Set("Content-Type", "application/json")
  // json.NewEncoder(w).Encode(response)

	rows, err := db.DB.Query("SELECT Task_ID, Task_Name, Task_Done FROM Tasks WHERE User_ID= ?",userID)
	if err != nil {
		http.Error(w, `{"error":"データベースエラー"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Done)
		if err != nil {
			http.Error(w,`{"error": "データ取得エラー"}`, http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

//Task追加処理
func AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// データベースにタスクを追加
	result, err := db.DB.Exec("INSERT INTO Tasks (Task_Name, Task_Done) VALUES (?, ?)", task.Name, task.Done)
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "ID取得エラー", http.StatusInternalServerError)
		return
	}
	task.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

//Task削除処理
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "無効なIDです", http.StatusBadRequest)
		return
	}

	// データベースからタスクを削除
	result, err := db.DB.Exec("DELETE FROM Tasks WHERE Task_Id = ?", id)
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "指定されたIDのタスクが見つかりません", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//Task状態変更処理
func ToggleTaskDone(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "無効なIDです", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	// グローバルな tasks スライスを参照する.
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("UPDATE Tasks SET Task_Done = ? WHERE Task_Id = ?", updatedTask.Done, id)
	if err != nil {
			http.Error(w, "データベースエラー", http.StatusInternalServerError)
			return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
			http.Error(w, "指定されたIDのタスクが見つかりません", http.StatusNotFound)
			return
	}

	// 更新されたタスクをクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)

}
