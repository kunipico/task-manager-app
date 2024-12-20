package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	// データベースに接続
	var err error
	dsn := "mysql:mysql#MYSQL123@tcp(task-manager-db.c0xkgufvdnai.ap-northeast-1.rds.amazonaws.com:3306)/TaskManager?charset=utf8mb4&parseTime=true"
	// dsn := "mysql:mysql#MYSQL123@tcp(db:3306)/TaskManager?charset=utf8mb4&parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("データベース接続エラー: ", err)
	}

	// DB接続確認
	// if err := DB.Ping(); err != nil {
	// 	log.Fatal("データベース接続エラー: ", err)
	// }else{
	// 	log.Fatal("データベース接続成功！！！")
	// }
}