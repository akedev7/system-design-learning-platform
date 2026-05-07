CREATE TABLE user_lesson_progress (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    score INTEGER NOT NULL DEFAULT 0,
    passed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, lesson_id)
);

CREATE INDEX idx_user_lesson_progress_user_id ON user_lesson_progress(user_id);
CREATE INDEX idx_user_lesson_progress_lesson_id ON user_lesson_progress(lesson_id);