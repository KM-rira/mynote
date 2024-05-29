CREATE DATABASE IF NOT EXISTS mynote_db CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE mynote_db;

CREATE TABLE IF NOT EXISTS note (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    contents TEXT NOT NULL,
    category VARCHAR(255),
    important BOOLEAN,
    created_at DATETIME,
    updated_at DATETIME
) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- サンプルデータの挿入
INSERT INTO note (title, contents, category, important, created_at, updated_at)
VALUES
    ('Sample Note 1', 'This is the content of sample note 1', 'General', TRUE, NOW(), NOW()),
    ('Sample Note 2', 'This is the content of sample note 2', 'Work', FALSE, NOW(), NOW()),
    ('Sample Note 3', 'This is the content of sample note 3', 'Personal', TRUE, NOW(), NOW()),
    ('サンプルデータ4です。', 'こんにちは', 'テスト', TRUE, NOW(), NOW());
