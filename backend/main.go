package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

//認証情報
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


// タスク構造体
type Task struct {
	ID   int    `json:"Task_Id"`
	Name string `json:"Task_Name"`
	Done bool   `json:"Task_Done"`
}

var db *sql.DB

func init() {
	// データベースに接続
	var err error
	dsn := "mysql:mysql#MYSQL123@tcp(db:3306)/TaskManager?charset=utf8mb4"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("データベース接続エラー: ", err)
	}

	// DB接続確認
	if err := db.Ping(); err != nil {
		log.Fatal("データベース接続エラー: ", err)
	}
}

func main() {
	mux := http.NewServeMux() //マルチプレクサ。HTTPメソッドを指定してハンドラを呼び分けられるように登録可能。

	mux.HandleFunc("/login", login)
	mux.HandleFunc("GET /tasks",getTasks)
	mux.HandleFunc("POST /tasks",addTask)
	mux.HandleFunc("DELETE /tasks",deleteTask)
	mux.HandleFunc("PUT /tasks",toggleTaskDone)

	// CORSミドルウェアを設定
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	// AllowCredentials: true,
	})

	// CORSミドルウェアをHTTPサーバーに適用
	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

//ログイン処理
func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
	}

	// 固定のユーザー名とパスワードで認証
	if creds.Username == "user" && creds.Password == "password" {
			// ログイン成功
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Login successful"))
	} else {
			// 認証失敗
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}


//Task一覧取得処理
func getTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT Task_Id, Task_Name, Task_Done FROM Tasks")
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Done)
		if err != nil {
			http.Error(w, "データ取得エラー", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

//Task追加処理
func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// データベースにタスクを追加
	result, err := db.Exec("INSERT INTO Tasks (Task_Name, Task_Done) VALUES (?, ?)", task.Name, task.Done)
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
func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "無効なIDです", http.StatusBadRequest)
		return
	}

	// データベースからタスクを削除
	result, err := db.Exec("DELETE FROM Tasks WHERE Task_Id = ?", id)
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
func toggleTaskDone(w http.ResponseWriter, r *http.Request) {
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

	result, err := db.Exec("UPDATE Tasks SET Task_Done = ? WHERE Task_Id = ?", updatedTask.Done, id)
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

