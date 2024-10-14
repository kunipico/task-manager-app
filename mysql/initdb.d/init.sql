#!/bin/bash

-- データベースの作成
CREATE DATABASE IF NOT EXISTS TaskManager;

USE TaskManager;

-- ユーザーの作成
CREATE USER IF NOT EXISTS 'mysql'@'%' IDENTIFIED BY 'mysql#MYSQL123';

-- 権限の付与
GRANT ALL PRIVILEGES ON TaskManager.* TO 'mysql'@'%';

-- 既存テーブルの削除
DROP TABLE IF EXISTS Tasks;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Docs;

-- Tasksテーブル
CREATE TABLE Tasks (
    Task_ID INT AUTO_INCREMENT PRIMARY KEY,
    Task_Name VARCHAR(255) NOT NULL,
    Task_Description TEXT,
    Task_Done ENUM('Standby','Inprogress','Done'),
    Create_At DATETIME
)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Usersテーブル
CREATE TABLE Users (
    User_ID INT AUTO_INCREMENT PRIMARY KEY,
    User_Name VARCHAR(255) NOT NULL,
    Password VARCHAR(255) NOT NULL
)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Docsテーブル
CREATE TABLE Docs (
    Docs_ID INT AUTO_INCREMENT PRIMARY KEY,
    Task_ID INT NOT NULL,
    Create_At DATETIME,
    FOREIGN KEY (Task_ID) REFERENCES Tasks(Task_ID) ON DELETE CASCADE
)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;


-- ALTER TABLE Tasks CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- サンプルデータの挿入
INSERT INTO Tasks VALUES
(1,'チーム開発','チーム開発の説明','Standby','2024-10-12 09:00:00'),
(2,'読書','確率思考の戦略論を読書中。学んだことを記録していく。','Standby','2024-10-13 09:00:00'),
(3,'テスト問題作成','テスト問題を作成したけどクオリティが低かった。いや酷かった。','Standby','2024-10-14 09:00:00')
