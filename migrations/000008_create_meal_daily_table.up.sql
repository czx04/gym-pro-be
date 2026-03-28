CREATE TABLE meal_daily (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    date DATE NOT NULL,
    target_calories DECIMAL(10, 2) NOT NULL DEFAULT 0,
    target_protein_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    target_carbs_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    target_fat_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_meal_daily_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_date UNIQUE(user_id, date)
);
