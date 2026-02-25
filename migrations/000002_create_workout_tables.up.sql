-- Exercises table (pre-populated library)
CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    muscle_groups JSONB,
    equipment_needed JSONB,
    difficulty_level VARCHAR(50) NOT NULL,
    calories_per_minute DECIMAL(5, 2),
    video_url TEXT,
    thumbnail_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_difficulty CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
    CONSTRAINT check_category CHECK (category IN ('cardio', 'strength', 'flexibility', 'stretching'))
);

CREATE INDEX idx_exercises_category ON exercises(category);
CREATE INDEX idx_exercises_difficulty ON exercises(difficulty_level);
CREATE INDEX idx_exercises_is_active ON exercises(is_active);
CREATE INDEX idx_exercises_name ON exercises(name);

-- Workout Plans table
CREATE TABLE IF NOT EXISTS workout_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    difficulty_level VARCHAR(50) NOT NULL,
    estimated_duration_mins INTEGER,
    estimated_calories INTEGER,
    is_template BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_plan_difficulty CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced'))
);

CREATE INDEX idx_workout_plans_user ON workout_plans(user_id);
CREATE INDEX idx_workout_plans_public ON workout_plans(is_public) WHERE is_public = TRUE;

-- Workout Plan Exercises (M2M with configuration)
CREATE TABLE IF NOT EXISTS workout_plan_exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_plan_id UUID NOT NULL,
    exercise_id UUID NOT NULL,
    "order" INTEGER NOT NULL,
    sets INTEGER,
    reps INTEGER,
    duration_secs INTEGER,
    rest_secs INTEGER,
    notes TEXT,
    
    FOREIGN KEY (workout_plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    CONSTRAINT unique_plan_exercise_order UNIQUE (workout_plan_id, "order")
);

CREATE INDEX idx_workout_plan_exercises_plan ON workout_plan_exercises(workout_plan_id);
CREATE INDEX idx_workout_plan_exercises_exercise ON workout_plan_exercises(exercise_id);

-- Workout Schedules table
CREATE TABLE IF NOT EXISTS workout_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_plan_id UUID NOT NULL,
    user_id UUID NOT NULL,
    scheduled_date DATE NOT NULL,
    scheduled_time TIME,
    recurrence_rule VARCHAR(100),
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (workout_plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_workout_schedules_user ON workout_schedules(user_id);
CREATE INDEX idx_workout_schedules_date ON workout_schedules(scheduled_date);
CREATE INDEX idx_workout_schedules_user_date ON workout_schedules(user_id, scheduled_date);
CREATE INDEX idx_workout_schedules_completed ON workout_schedules(is_completed);

-- Workout Sessions table (actual workout tracking)
CREATE TABLE IF NOT EXISTS workout_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_schedule_id UUID,
    user_id UUID NOT NULL,
    workout_plan_id UUID NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_mins INTEGER,
    total_calories_burned INTEGER,
    notes TEXT,
    mood VARCHAR(50),
    difficulty_rating INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (workout_schedule_id) REFERENCES workout_schedules(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (workout_plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
    CONSTRAINT check_mood CHECK (mood IN ('happy', 'neutral', 'tired', 'energetic') OR mood IS NULL),
    CONSTRAINT check_difficulty_rating CHECK (difficulty_rating IS NULL OR (difficulty_rating >= 1 AND difficulty_rating <= 5))
);

CREATE INDEX idx_workout_sessions_user ON workout_sessions(user_id);
CREATE INDEX idx_workout_sessions_started ON workout_sessions(started_at);
CREATE INDEX idx_workout_sessions_user_started ON workout_sessions(user_id, started_at);

-- Workout Session Exercises table (per-exercise tracking)
CREATE TABLE IF NOT EXISTS workout_session_exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_session_id UUID NOT NULL,
    exercise_id UUID NOT NULL,
    "order" INTEGER NOT NULL,
    target_sets INTEGER,
    target_reps INTEGER,
    actual_sets_completed JSONB,
    duration_secs INTEGER,
    notes TEXT,
    skipped BOOLEAN DEFAULT FALSE,
    
    FOREIGN KEY (workout_session_id) REFERENCES workout_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

CREATE INDEX idx_workout_session_exercises_session ON workout_session_exercises(workout_session_id);
CREATE INDEX idx_workout_session_exercises_exercise ON workout_session_exercises(exercise_id);
