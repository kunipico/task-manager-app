package task

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"task-manager-api/db"

	"github.com/golang-jwt/jwt"
)

// タスク構造体
type Task struct {
	ID   int    `json:"Task_ID"`
	Name string `json:"Task_Name"`
  Details string `json:"Task_Details"`
	Done string   `json:"Task_Done"`
}

// JWTシークレットキー
var jwtKey = []byte("my_secret_key")

// JWTのペイロードを定義
type Claims struct {
	Userinfo string `json:"userinfo"`
	jwt.StandardClaims
}

// タイムスタンプ生成
func MakeTimeStamp()(time.Time, error){
  // 日本標準時をロード
  jst, err := time.LoadLocation("Asia/Tokyo")
  if err != nil {
    return time.Time{},fmt.Errorf("time zone error:%w",err)
  }
  // 現在の日本時間を取得
  currentTime := time.Now().In(jst)
  fmt.Println("currentTime : ",currentTime)
	return currentTime,nil
}



// CookieとJWTの検証処理
func ExtractUserIDAndParam(r *http.Request) (string, string, error) {
	// Cookieの内容を確認
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Nothing Cookie")
		return "", "", fmt.Errorf("nothing Cookie")
	}
	// tokenStr := cookie.Value
	// JWTの検証
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("Unauthorized")
	}

	// クレームからユーザーIDを取得
	userID := claims.Userinfo
	
	// クレーム内のIDをDBから検索
	err = db.DB.QueryRow("SELECT User_ID FROM Users WHERE User_ID = ?", userID).Scan()
	if err == sql.ErrNoRows {
		return "", "", fmt.Errorf("Unauthorized")
	}

	// URLパスからパラメータを取得
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	return userID, idStr, nil
}



//Task一覧取得処理
func GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, _, err := ExtractUserIDAndParam(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query("SELECT Task_ID, Task_Name, Task_Details, Task_Done FROM Tasks WHERE User_ID= ?",userID)
	if err != nil {
		http.Error(w, `{"error":"データベースエラー"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name,&task.Details, &task.Done)
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
	userID, _, err := ExtractUserIDAndParam(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}

  var task Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
    fmt.Print("JSONデコードエラー")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentTime,err := MakeTimeStamp()
	if err != nil {
		fmt.Println("Error generating timestamp:", err)
		return
	}

	// データベースにタスクを追加
	_, err = db.DB.Exec("INSERT INTO Tasks (User_ID, Task_Name, Task_Details, Task_Done, Create_At) VALUES (?, ?, ?, ?, ?)", userID, task.Name,task.Details, task.Done, currentTime)
	if err != nil {
		http.Error(w, `{"error":"データベースエラー"}`, http.StatusInternalServerError)
		return
	}
}



//Task削除処理
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Cookieの内容を確認
	_, idStr, err := ExtractUserIDAndParam(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "無効なIDです", http.StatusBadRequest)
		return
	}

	fmt.Println("idStr, id : ",idStr,id)

	// タスク状態取得
	var done string
	err = db.DB.QueryRow("SELECT Task_Done FROM Tasks WHERE Task_ID = ?", id).Scan(&done)
	if err != nil {
		http.Error(w, `{"error":"データベースエラー1"}`, http.StatusInternalServerError)
		return
	}

	fmt.Println("done: ",done)

	if done == "Done" {
		// データベースからタスクを削除
		_, err := db.DB.Exec("UPDATE Tasks SET Task_Done = 'Standby' WHERE Task_ID = ?", id)
		if err != nil {
			http.Error(w, "データベースエラー", http.StatusInternalServerError)
			return
		}
	}else{
		_, err := db.DB.Exec("UPDATE Tasks SET Task_Done = 'Done' WHERE Task_ID = ?", id)
		if err != nil {
			http.Error(w, "データベースエラー", http.StatusInternalServerError)
			return
		}
		if done == "Inprogress"{
			currentTime,err := MakeTimeStamp()
			if err != nil {
				fmt.Println("Error generating timestamp:", err)
				return
			}
			// 実行中→停止　時間記録
			_, err = db.DB.Exec("INSERT INTO Times (Task_ID, SetStatus,SetTime) VALUES (?,?,?)",id,"Stop",currentTime)
			if err != nil {
				http.Error(w, `{"error":"データベースエラー3"}`, http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
	
}



//Task状態変更処理
func ToggleTaskDone(w http.ResponseWriter, r *http.Request) {
	currentTime,err := MakeTimeStamp()
	if err != nil {
		fmt.Println("Error generating timestamp:", err)
		return
	}
	
	userID, idStr, err := ExtractUserIDAndParam(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "無効なIDです", http.StatusBadRequest)
		return
	}

	fmt.Println("idStr, id : ",idStr,id)
	
	var updatedTask Task
	// グローバルな tasks スライスを参照する.
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// status初期化処理
	// 実行中のタスク検索
	var InprogressTaskId int
	rows, err := db.DB.Query("SELECT Task_ID FROM Tasks WHERE User_ID = ? AND Task_Done = 'Inprogress'",userID)
	if err != nil {
		http.Error(w, `{"error":"データベースエラー1"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&InprogressTaskId)
		if err != nil {
			http.Error(w,`{"error": "データ取得エラー"}`, http.StatusInternalServerError)
			return
		}
		// 実行中→停止　状態変更
		_, err = db.DB.Exec("UPDATE Tasks SET Task_Done = 'Standby' WHERE Task_ID = ?",InprogressTaskId)
		if err != nil {
			http.Error(w, `{"error":"データベースエラー2"}`, http.StatusInternalServerError)
			return
		}
		// 実行中→停止　時間記録
		_, err = db.DB.Exec("INSERT INTO Times (Task_ID, SetStatus,SetTime) VALUES (?,?,?)",InprogressTaskId,"Stop",currentTime)
		if err != nil {
			http.Error(w, `{"error":"データベースエラー3"}`, http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("updatedTask.Done, id: ", updatedTask.Done, id)

	if updatedTask.Done == "Inprogress"{
		result, err := db.DB.Exec("UPDATE Tasks SET Task_Done = ? WHERE Task_ID = ?", updatedTask.Done, id)
		if err != nil {
				http.Error(w, `{"error":"データベースエラー4"}`, http.StatusInternalServerError)
				return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
				http.Error(w, `{"error":"指定されたIDのタスクが見つかりません"}`, http.StatusNotFound)
				return
		}
		// 停止→実行中　時間記録
		_, err = db.DB.Exec("INSERT INTO Times (Task_ID, SetStatus,SetTime) VALUES (?,?,?)",id,"Start",currentTime)
		if err != nil {
			http.Error(w, `{"error":"データベースエラー5"}`, http.StatusInternalServerError)
			return
		}
	}

	// 更新されたタスクをクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)

}