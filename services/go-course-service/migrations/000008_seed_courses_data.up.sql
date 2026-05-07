-- Seed System Design Fundamentals course
-- Migration 000008_seed_courses_data.up.sql

-- Course 1: System Design Fundamentals
INSERT INTO courses (title, description) VALUES 
('System Design Fundamentals', 'Learn the core principles of designing scalable systems.');

-- Modules for System Design Fundamentals
WITH course1 AS (
    SELECT id FROM courses WHERE title = 'System Design Fundamentals'
)
INSERT INTO modules (course_id, title, description, order_index)
SELECT course1.id, 'Introduction to System Design', 'Basic concepts and terminology', 1
FROM course1;

INSERT INTO modules (course_id, title, description, order_index)
SELECT course1.id, 'Client-Server Architecture', 'Understanding client-server model', 2
FROM course1;

-- Lessons for Module 1 (Introduction to System Design)
WITH m1 AS (
    SELECT id FROM modules WHERE title = 'Introduction to System Design'
)
INSERT INTO lessons (module_id, title, description, order_index)
SELECT m1.id, 'What is System Design?', 'Introduction to system design principles', 1
FROM m1;

WITH m1 AS (
    SELECT id FROM modules WHERE title = 'Introduction to System Design'
)
INSERT INTO lessons (module_id, title, description, order_index)
SELECT m1.id, 'Key Components', 'Building blocks of systems', 2
FROM m1;

-- Course 2: Scalability & Load Balancing
INSERT INTO courses (title, description) VALUES 
('Scalability & Load Balancing', 'Master horizontal scaling and load distribution.');

-- Course 3: Database Choices
INSERT INTO courses (title, description) VALUES 
('Database Choices', 'Select the right database for your use case.');

-- Modules for Scalability & Load Balancing
WITH course2 AS (
    SELECT id FROM courses WHERE title = 'Scalability & Load Balancing'
)
INSERT INTO modules (course_id, title, description, order_index)
SELECT course2.id, 'Horizontal vs Vertical Scaling', 'Understanding scaling strategies', 1
FROM course2;

-- Modules for Database Choices
WITH course3 AS (
    SELECT id FROM courses WHERE title = 'Database Choices'
)
INSERT INTO modules (course_id, title, description, order_index)
SELECT course3.id, 'SQL vs NoSQL', 'Choosing the right database type', 1
FROM course3;