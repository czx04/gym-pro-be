CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    oauth_provider VARCHAR(50),
    oauth_id VARCHAR(255),
    name VARCHAR(100) NOT NULL,
    bio TEXT,
    avatar_url TEXT,
    date_of_birth DATE,
    gender VARCHAR(20),
    height_cm DECIMAL(5, 2),
    weight_kg DECIMAL(5, 2),
    fitness_goal VARCHAR(50),
    activity_level VARCHAR(50),
    daily_calorie_target INTEGER,
    protein_target_g INTEGER,
    carbs_target_g INTEGER,
    fat_target_g INTEGER,
    privacy_settings JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_oauth UNIQUE (oauth_provider, oauth_id),
    CONSTRAINT check_height CHECK (height_cm IS NULL OR (height_cm > 0 AND height_cm <= 300)),
    CONSTRAINT check_weight CHECK (weight_kg IS NULL OR (weight_kg > 0 AND weight_kg <= 500)),
    CONSTRAINT check_calories CHECK (daily_calorie_target IS NULL OR (daily_calorie_target >= 500 AND daily_calorie_target <= 10000))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_id) WHERE oauth_provider IS NOT NULL;
