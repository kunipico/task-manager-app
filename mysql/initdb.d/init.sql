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

-- 論文テーブル
CREATE TABLE Tasks (
    Task_ID INT AUTO_INCREMENT PRIMARY KEY,
    Task_Name VARCHAR(255) NOT NULL,
    Task_Done BOOLEAN NOT NULL
);

ALTER TABLE Tasks CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- サンプルデータの挿入
INSERT INTO Tasks VALUES
(1,"チーム開発",False),
(2,"読書",True),
(3,"テスト問題作成",True)
