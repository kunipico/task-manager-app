package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// タスク構造体
type Task struct {
	ID   int    `json:"Task_id"`
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

// var (
// 	tasks   []Task
// 	nextID  int
// )

// func init() {
// 	tasks = []Task{
// 		{ID: 1, Name: "インターン研修", Done: false},
// 		{ID: 2, Name: "筋トレ", Done: false},
// 		{ID: 3, Name: "休息", Done: false},
// 	}
// 	nextID = 4
// }

func main() {
	//リクエストハンドラ
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", tasksHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))

}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
			getTasks(w, r) // 全てのtaskの取得
	case "POST":
			addTask(w, r) // 新しいtaskの追加
	case "DELETE":
			deleteTask(w, r) // 指定のtaskを削除
	default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

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

func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// データベースにタスクを追加
	result, err := db.Exec("INSERT INTO Tasks (name, done) VALUES (?, ?)", task.Name, task.Done)
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

// func getTasks(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(tasks)
// }

// func addTask(w http.ResponseWriter, r *http.Request) {
// 	var task Task
// 	err := json.NewDecoder(r.Body).Decode(&task)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	task.ID = nextID
// 	nextID++
// 	tasks = append(tasks, task)

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(task)
// }

// func deleteTask(w http.ResponseWriter, r *http.Request) {
// 	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "無効なIDです", http.StatusBadRequest)
// 		return
// 	}

// 	for i, task := range tasks {
// 		if task.ID == id {
// 			tasks = append(tasks[:i], tasks[i+1:]...)
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 	}
// 	http.Error(w, "指定されたIDのタスクが見つかりません", http.StatusNotFound)
// }