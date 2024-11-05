package main

import (
	"fmt"
	"log"
	"net/http"

	"task-manager-api/auth"
	"task-manager-api/db"
	"task-manager-api/task"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

func main() {
	db.InitDB()
	fmt.Println("Hello golang")
	mux := http.NewServeMux() //マルチプレクサ。HTTPメソッドを指定してハンドラを呼び分けられるように登録可能。

	mux.HandleFunc("/login", auth.Login)
  mux.HandleFunc("/logout", auth.Logout)
  mux.HandleFunc("POST /signup", auth.Signup)
  mux.HandleFunc("/users", auth.Users)    //ユーザ登録確認用
	mux.HandleFunc("GET /tasks",task.GetTasks)
	mux.HandleFunc("POST /tasks",task.AddTask)
	mux.HandleFunc("DELETE /tasks",task.DeleteTask)
	mux.HandleFunc("PUT /tasks",task.ToggleTaskDone)

	// CORSミドルウェアを設定
	c := cors.New(cors.Options{
		// AllowedOrigins: []string{"http://localhost:3000/","http://localhost:3000/login"},
    AllowedOrigins: []string{"http://localhost:3000","http://next-app:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE","OPTION"},
		AllowedHeaders: []string{"Content-Type","Authorization"},
	  AllowCredentials: true,

	})

	// CORSミドルウェアをHTTPサーバーに適用
	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

