-- Remove seeded course data
-- Migration 000008_seed_courses_data.down.sql

DELETE FROM lessons 
WHERE module_id IN (
    SELECT id FROM modules 
    WHERE course_id IN (
        SELECT id FROM courses WHERE title IN (
            'System Design Fundamentals',
            'Scalability & Load Balancing', 
            'Database Choices'
        )
    )
);

DELETE FROM modules 
WHERE course_id IN (
    SELECT id FROM courses WHERE title IN (
        'System Design Fundamentals',
        'Scalability & Load Balancing',
        'Database Choices'
    )
);

DELETE FROM courses WHERE title IN (
    'System Design Fundamentals',
    'Scalability & Load Balancing',
    'Database Choices'
);