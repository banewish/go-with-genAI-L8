-- Seed demo data
INSERT INTO categories (name, created_at, updated_at) VALUES
('Fiction', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Science', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('History', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO authors (full_name, bio, created_at, updated_at) VALUES
('George R. R. Martin', 'American novelist', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('J. R. R. Tolkien', 'British writer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Stephen King', 'Horror novelist', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Haruki Murakami', 'Japanese writer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Agatha Christie', 'Mystery writer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO books (title, isbn, category_id, description, created_at, updated_at) VALUES
('A Game of Thrones', '978-0553103540', 1, 'First book in ASOIAF series', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('The Fellowship of the Ring', '978-0547928227', 1, 'First book of LOTR trilogy', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('The Shining', '978-0385333312', 2, 'Psychological horror novel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Norwegian Wood', '978-0679442882', 1, 'Coming-of-age novel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Murder on the Orient Express', '978-0062073556', 3, 'Mystery novel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Dune', '978-0441173570', 2, 'Science fiction epic', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Foundation', '978-0553293357', 2, 'Science fiction masterpiece', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('The Catcher in the Rye', '978-0316769174', 1, 'Coming-of-age classic', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('1984', '978-0451524935', 3, 'Dystopian novel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Brave New World', '978-0060085345', 3, 'Science fiction classic', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO book_authors (book_id, author_id) VALUES
(1, 1), (2, 2), (3, 3), (4, 4), (5, 5), (6, 1), (7, 2), (8, 3), (9, 4), (10, 5);
