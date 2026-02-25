-- Foods table (library - system + user custom)
CREATE TABLE IF NOT EXISTS foods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    image_url TEXT,
    description TEXT,
    brand VARCHAR(100),
    serving_size DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    calories DECIMAL(8, 2) NOT NULL,
    protein_g DECIMAL(7, 2) NOT NULL,
    carbs_g DECIMAL(7, 2) NOT NULL,
    fat_g DECIMAL(7, 2) NOT NULL,
    fiber_g DECIMAL(7, 2),
    is_system BOOLEAN DEFAULT FALSE,
    created_by_user_id UUID,
    category VARCHAR(50),
    barcode VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_serving_size CHECK (serving_size > 0),
    CONSTRAINT check_calories CHECK (calories >= 0),
    CONSTRAINT check_protein CHECK (protein_g >= 0),
    CONSTRAINT check_carbs CHECK (carbs_g >= 0),
    CONSTRAINT check_fat CHECK (fat_g >= 0),
    CONSTRAINT check_category CHECK (category IN ('protein', 'carb', 'vegetable', 'fruit', 'dairy', 'fat', 'snack', 'beverage', 'other') OR category IS NULL)
);

CREATE INDEX idx_foods_name ON foods(name);
CREATE INDEX idx_foods_category ON foods(category);
CREATE INDEX idx_foods_is_system ON foods(is_system);
CREATE INDEX idx_foods_user ON foods(created_by_user_id) WHERE created_by_user_id IS NOT NULL;
CREATE INDEX idx_foods_barcode ON foods(barcode) WHERE barcode IS NOT NULL;

-- Recipes table
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    prep_time_mins INTEGER,
    cook_time_mins INTEGER,
    servings INTEGER NOT NULL,
    instructions TEXT,
    image_url TEXT,
    total_calories DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_protein_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    total_carbs_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    total_fat_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    per_serving_calories DECIMAL(10, 2) NOT NULL DEFAULT 0,
    per_serving_protein_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    per_serving_carbs_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    per_serving_fat_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    is_public BOOLEAN DEFAULT FALSE,
    visibility VARCHAR(20) NOT NULL DEFAULT 'private',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_servings CHECK (servings > 0),
    CONSTRAINT check_visibility CHECK (visibility IN ('public', 'private', 'friends'))
);

CREATE INDEX idx_recipes_user ON recipes(user_id);
CREATE INDEX idx_recipes_public ON recipes(is_public) WHERE is_public = TRUE;
CREATE INDEX idx_recipes_name ON recipes(name);

-- Recipe Foods (M2M with quantities)
CREATE TABLE IF NOT EXISTS recipe_foods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL,
    food_id UUID NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    calories DECIMAL(8, 2) NOT NULL,
    protein_g DECIMAL(7, 2) NOT NULL,
    carbs_g DECIMAL(7, 2) NOT NULL,
    fat_g DECIMAL(7, 2) NOT NULL,
    
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    CONSTRAINT check_rf_quantity CHECK (quantity > 0)
);

CREATE INDEX idx_recipe_foods_recipe ON recipe_foods(recipe_id);
CREATE INDEX idx_recipe_foods_food ON recipe_foods(food_id);

-- Meal Logs table (actual consumption tracking)
CREATE TABLE IF NOT EXISTS meal_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    log_date DATE NOT NULL,
    meal_time VARCHAR(20) NOT NULL,
    total_calories DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_protein_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    total_carbs_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    total_fat_g DECIMAL(8, 2) NOT NULL DEFAULT 0,
    notes TEXT,
    mood VARCHAR(20),
    energy_level INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_meal_time CHECK (meal_time IN ('breakfast', 'lunch', 'dinner', 'snack', 'other')),
    CONSTRAINT check_mood CHECK (mood IN ('great', 'good', 'okay', 'tired') OR mood IS NULL),
    CONSTRAINT check_energy CHECK (energy_level IS NULL OR (energy_level >= 1 AND energy_level <= 5))
);

CREATE INDEX idx_meal_logs_user ON meal_logs(user_id);
CREATE INDEX idx_meal_logs_date ON meal_logs(log_date);
CREATE INDEX idx_meal_logs_user_date ON meal_logs(user_id, log_date);
CREATE INDEX idx_meal_logs_user_date_time ON meal_logs(user_id, log_date, meal_time);

-- Meal Log Items table
CREATE TABLE IF NOT EXISTS meal_log_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meal_log_id UUID NOT NULL,
    item_type VARCHAR(20) NOT NULL,
    food_id UUID,
    recipe_id UUID,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    serving_size DECIMAL(5, 2),
    calories DECIMAL(8, 2) NOT NULL,
    protein_g DECIMAL(7, 2) NOT NULL,
    carbs_g DECIMAL(7, 2) NOT NULL,
    fat_g DECIMAL(7, 2) NOT NULL,
    "order" INTEGER NOT NULL,
    
    FOREIGN KEY (meal_log_id) REFERENCES meal_logs(id) ON DELETE CASCADE,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    CONSTRAINT check_item_type CHECK (item_type IN ('food', 'recipe')),
    CONSTRAINT check_item_reference CHECK (
        (item_type = 'food' AND food_id IS NOT NULL AND recipe_id IS NULL) OR
        (item_type = 'recipe' AND recipe_id IS NOT NULL AND food_id IS NULL)
    ),
    CONSTRAINT check_mli_quantity CHECK (quantity > 0)
);

CREATE INDEX idx_meal_log_items_log ON meal_log_items(meal_log_id);
CREATE INDEX idx_meal_log_items_food ON meal_log_items(food_id) WHERE food_id IS NOT NULL;
CREATE INDEX idx_meal_log_items_recipe ON meal_log_items(recipe_id) WHERE recipe_id IS NOT NULL;
