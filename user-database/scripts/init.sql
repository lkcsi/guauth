CREATE DATABASE user_db;
use user_db;

CREATE TABLE users (
	username VARCHAR(255) PRIMARY KEY,
	password VARCHAR(255) NOT NULL
);

INSERT INTO users (username, password)
VALUES 
	('guest', '$2a$10$KJKTiTcOhHjVIVH74u8pCOv18tzOs4Fd8bd8Dl7mZlJy/q2Tj2Vjq')
