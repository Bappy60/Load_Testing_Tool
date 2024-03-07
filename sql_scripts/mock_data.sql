-- mock_data.sql
CREATE DATABASE IF NOT EXISTS demodb;
USE demodb;
INSERT INTO books (name, publication_year, number_of_pages, author_id, publication) VALUES
('Book 1', 2020, 300, 1, 'Publisher 1'),
('Book 2', 2018, 250, 2, 'Publisher 2'),
('Book 3', 2019, 280, 3, 'Publisher 3'),
('Book 4', 2021, 320, 4, 'Publisher 4'),
('Book 5', 2017, 270, 5, 'Publisher 5'),
('Book 6', 2022, 350, 6, 'Publisher 6'),
('Book 7', 2016, 230, 7, 'Publisher 7'),
('Book 8', 2015, 200, 8, 'Publisher 8'),
('Book 9', 2014, 180, 9, 'Publisher 9'),
('Book 10', 2013, 150, 10, 'Publisher 10');

INSERT INTO authors (name, email, age) VALUES
('Author 1', 'author1@example.com', 30),
('Author 2', 'author2@example.com', 35),
('Author 3', 'author3@example.com', 40),
('Author 4', 'author4@example.com', 45),
('Author 5', 'author5@example.com', 50),
('Author 6', 'author6@example.com', 55),
('Author 7', 'author7@example.com', 60),
('Author 8', 'author8@example.com', 65),
('Author 9', 'author9@example.com', 70),
('Author 10', 'author10@example.com', 75);
