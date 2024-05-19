CREATE DATABASE IF NOT EXISTS mynote_db;

USE mynote_db;

CREATE TABLE IF NOT EXISTS note (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(255),
    important_flag BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- sample data
INSERT INTO note (title, content, category, important_flag) VALUES
    ('Meeting Notes', 'Discuss project timeline and deliverables.', 'Work', TRUE),
    ('Shopping List', 'Milk, Bread, Eggs, and Butter.', 'Personal', FALSE),
    ('Reading List', 'Read "The Pragmatic Programmer".', 'Personal', TRUE),
    ('Workout Plan', 'Cardio on Monday, Strength training on Wednesday.', 'Health', FALSE);
