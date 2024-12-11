package main

import (
	"fmt"
	"log"
	"net/http"

	"task-manager-api/auth"
	"task-manager-api/db"
	"task-manager-api/task"
	"task-manager-api/times"

	"task-manager-api/docs"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

func main() {
	db.InitDB()
	fmt.Println("Hello golang")
	mux := http.NewServeMux() //マルチプレクサ。HTTPメソッドを指定してハンドラを呼び分けられるように登録可能。

	mux.HandleFunc("POST /login", auth.Login)
  mux.HandleFunc("DELETE /logout", auth.Logout)
  mux.HandleFunc("POST /signup", auth.Signup)
  // mux.HandleFunc("/users", auth.Users)    //ユーザ登録確認用
	mux.HandleFunc("GET /tasks",task.GetTasks)
	mux.HandleFunc("POST /tasks",task.AddTask)
	mux.HandleFunc("DELETE /tasks/{id}",task.DeleteTask)
	mux.HandleFunc("PUT /tasks/{id}",task.ToggleTaskDone)
	mux.HandleFunc("GET /tasks/time-info/{id}",times.GetTimeInfo)
	mux.HandleFunc("GET /tasks/documents/{id}",docs.GetDocuments)
	mux.HandleFunc("POST /tasks/documents/{id}",docs.AddDocument)


	// CORSミドルウェアを設定
	c := cors.New(cors.Options{
    AllowedOrigins: []string{"http://localhost:3000","http://next-app:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE","OPTIONS"},
		AllowedHeaders: []string{"Content-Type","Authorization"},
	  AllowCredentials: true,

	})

	// CORSミドルウェアをHTTPサーバーに適用
	handler := c.Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

