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

-- Usersテーブル
CREATE TABLE Users (
  User_ID INT AUTO_INCREMENT PRIMARY KEY,
  User_Name VARCHAR(255) NOT NULL,
  Emailaddress VARCHAR(50) NOT NULL,
  Password VARCHAR(255) NOT NULL
)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Tasksテーブル
CREATE TABLE Tasks (
  Task_ID INT AUTO_INCREMENT PRIMARY KEY,
  User_ID INT NOT NULL,
  Task_Name VARCHAR(255) NOT NULL,
  Task_Details TEXT,
  Task_Done ENUM('Standby','Inprogress','Done'),
  Create_At DATETIME,
  FOREIGN KEY (User_ID) REFERENCES Users(User_ID) ON DELETE CASCADE
)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Docsテーブル
CREATE TABLE Docs (
  Docs_ID INT AUTO_INCREMENT PRIMARY KEY,
  Task_ID INT NOT NULL,
  Documents TEXT NOT NULL,
  Create_At DATETIME,
  FOREIGN KEY (Task_ID) REFERENCES Tasks(Task_ID) ON DELETE CASCADE
)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Timesテーブル
CREATE TABLE Times (
  Time_ID INT AUTO_INCREMENT PRIMARY KEY,
  Task_ID INT NOT NULL,
  setStatus ENUM('Start','Stop'),
  setTime DATETIME,
  FOREIGN KEY (Task_ID) REFERENCES Tasks(Task_ID) ON DELETE CASCADE
)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- ALTER TABLE Tasks CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

INSERT INTO Users VALUES
(1,'morisato','morisato@example.com','$2a$10$Wq3R3BU9loIkuLMIP6sEFe4PMspbPMaQQcjbMfcJP4gQUZRO6tONK'),
(2,'kunihiko','kunihiko@example.com','$2a$10$vIjaaJYIq./6EDMlnZUQeucpNkPHxChWwLp6o85IKcvzgjbHsWWRe');

-- サンプルデータの挿入
INSERT INTO Tasks VALUES
(1,1,'個人開発','チーム開発の説明','Standby','2024-10-12 09:00:00'),
(2,1,'読書','確率思考の戦略論を読書中。学んだことを記録していく。','Standby','2024-10-13 09:00:00'),
(3,1,'テスト問題作成','テスト問題を作成したけどクオリティが低かった。いや酷かった。','Standby','2024-10-14 09:00:00'),
(4,2,'ランニング','目標タイム１５分','Standby','2024-10-12 09:00:00'),
(5,2,'囲碁','４段で勝ち越し！！！','Standby','2024-10-13 09:00:00'),
(6,2,'予定表作成','来年の予定を考える。','Standby','2024-10-14 09:00:00');
