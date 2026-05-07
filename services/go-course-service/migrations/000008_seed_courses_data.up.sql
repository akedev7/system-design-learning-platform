-- Seed System Design Fundamentals course
-- Migration 000008_seed_courses_data.up.sql

-- Course 1: System Design Fundamentals
INSERT INTO courses (title, description) VALUES 
('System Design Fundamentals', 'Learn the core principles of designing scalable systems.');

-- Modules for System Design Fundamentals
INSERT INTO modules (course_id, title, description, order_index)
SELECT id, 'Introduction to System Design', 'Basic concepts and terminology', 1
FROM courses WHERE title = 'System Design Fundamentals';

INSERT INTO modules (course_id, title, description, order_index)
SELECT id, 'Client-Server Architecture', 'Understanding client-server model', 2
FROM courses WHERE title = 'System Design Fundamentals';

-- Lessons for Module 1 (Introduction to System Design)
INSERT INTO lessons (module_id, title, description, order_index, content_jsonb)
SELECT id, 'What is System Design?', 'Introduction to system design principles', 1, 
  '[
    {"type": "Text", "order": 1, "config": {"content": "System design is the process of defining the architecture, components, modules, interfaces, and data for a system to satisfy specified requirements."}},
    {"type": "Quiz", "order": 2, "config": {"questions": [{"id": 1, "type": "multiple_choice", "question": "What is system design?", "options": ["Writing code", "Defining system architecture", "Testing software", "Deploying applications"], "correct": "Defining system architecture", "points": 10}]}},
    {"type": "ReactFlowDiagram", "order": 3, "config": {"nodeTypes": {"Client": 1, "Server": 1}, "edges": [{"from": "Client", "to": "Server"}]}}
  ]'::jsonb
FROM modules WHERE title = 'Introduction to System Design';

INSERT INTO lessons (module_id, title, description, order_index, content_jsonb)
SELECT id, 'Key Components', 'Building blocks of systems', 2,
  '[
    {"type": "Text", "order": 1, "config": {"content": "Key components include: Load Balancers, Servers, Databases, Caches, and CDNs."}},
    {"type": "CodeSnippet", "order": 2, "config": {"language": "javascript", "code": "const server = require('http').createServer();"}}
  ]'::jsonb
FROM modules WHERE title = 'Introduction to System Design';

-- Course 2: Scalability & Load Balancing
INSERT INTO courses (title, description) VALUES 
('Scalability & Load Balancing', 'Master horizontal scaling and load distribution.');

-- Course 3: Database Choices
INSERT INTO courses (title, description) VALUES 
('Database Choices', 'Select the right database for your use case.');

-- Modules for Scalability & Load Balancing
INSERT INTO modules (course_id, title, description, order_index)
SELECT id, 'Horizontal vs Vertical Scaling', 'Understanding scaling strategies', 1
FROM courses WHERE title = 'Scalability & Load Balancing';

-- Modules for Database Choices
INSERT INTO modules (course_id, title, description, order_index)
SELECT id, 'SQL vs NoSQL', 'Choosing the right database type', 1
FROM courses WHERE title = 'Database Choices';