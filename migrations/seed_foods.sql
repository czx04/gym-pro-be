-- Sample Food Data
-- Run this manually after migrations to populate food library
-- Usage: psql -U gymadmin -d gym_pro_db -f migrations/seed_foods.sql

-- Protein Foods
INSERT INTO foods (id, name, description, serving_size, unit, calories, protein_g, carbs_g, fat_g, image_url, is_system, category) VALUES
(uuid_generate_v4(), 'Chicken Breast', 'Boneless skinless chicken breast', 100, 'g', 165, 31, 0, 3.6, 'https://placehold.co/400x400?text=Chicken+Breast', true, 'protein'),
(uuid_generate_v4(), 'Salmon', 'Fresh atlantic salmon', 100, 'g', 206, 22, 0, 13, 'https://placehold.co/400x400?text=Salmon', true, 'protein'),
(uuid_generate_v4(), 'Eggs', 'Large whole eggs', 50, 'g', 72, 6.3, 0.4, 4.8, 'https://placehold.co/400x400?text=Eggs', true, 'protein'),
(uuid_generate_v4(), 'Greek Yogurt', 'Plain non-fat greek yogurt', 100, 'g', 59, 10, 3.6, 0.4, 'https://placehold.co/400x400?text=Greek+Yogurt', true, 'protein'),
(uuid_generate_v4(), 'Tofu', 'Firm tofu', 100, 'g', 76, 8, 1.9, 4.8, 'https://placehold.co/400x400?text=Tofu', true, 'protein');

-- Carbohydrate Foods
INSERT INTO foods (id, name, description, serving_size, unit, calories, protein_g, carbs_g, fat_g, fiber_g, image_url, is_system, category) VALUES
(uuid_generate_v4(), 'Brown Rice', 'Cooked brown rice', 100, 'g', 123, 2.6, 25.6, 1, 1.6, 'https://placehold.co/400x400?text=Brown+Rice', true, 'carb'),
(uuid_generate_v4(), 'Oatmeal', 'Rolled oats cooked', 100, 'g', 71, 2.5, 12, 1.5, 1.7, 'https://placehold.co/400x400?text=Oatmeal', true, 'carb'),
(uuid_generate_v4(), 'Sweet Potato', 'Baked sweet potato', 100, 'g', 90, 2, 21, 0.2, 3.3, 'https://placehold.co/400x400?text=Sweet+Potato', true, 'carb'),
(uuid_generate_v4(), 'Whole Wheat Bread', 'Whole wheat bread slice', 30, 'g', 80, 4, 14, 1, 2, 'https://placehold.co/400x400?text=Whole+Wheat+Bread', true, 'carb'),
(uuid_generate_v4(), 'Quinoa', 'Cooked quinoa', 100, 'g', 120, 4.4, 21.3, 1.9, 2.8, 'https://placehold.co/400x400?text=Quinoa', true, 'carb');

-- Vegetables
INSERT INTO foods (id, name, description, serving_size, unit, calories, protein_g, carbs_g, fat_g, fiber_g, image_url, is_system, category) VALUES
(uuid_generate_v4(), 'Broccoli', 'Raw broccoli', 100, 'g', 34, 2.8, 6.6, 0.4, 2.6, 'https://placehold.co/400x400?text=Broccoli', true, 'vegetable'),
(uuid_generate_v4(), 'Spinach', 'Raw spinach', 100, 'g', 23, 2.9, 3.6, 0.4, 2.2, 'https://placehold.co/400x400?text=Spinach', true, 'vegetable'),
(uuid_generate_v4(), 'Carrots', 'Raw carrots', 100, 'g', 41, 0.9, 9.6, 0.2, 2.8, 'https://placehold.co/400x400?text=Carrots', true, 'vegetable'),
(uuid_generate_v4(), 'Tomato', 'Fresh tomato', 100, 'g', 18, 0.9, 3.9, 0.2, 1.2, 'https://placehold.co/400x400?text=Tomato', true, 'vegetable'),
(uuid_generate_v4(), 'Bell Pepper', 'Red bell pepper', 100, 'g', 31, 1, 6, 0.3, 2.1, 'https://placehold.co/400x400?text=Bell+Pepper', true, 'vegetable');

-- Fruits
INSERT INTO foods (id, name, description, serving_size, unit, calories, protein_g, carbs_g, fat_g, fiber_g, image_url, is_system, category) VALUES
(uuid_generate_v4(), 'Banana', 'Medium banana', 118, 'g', 105, 1.3, 27, 0.4, 3.1, 'https://placehold.co/400x400?text=Banana', true, 'fruit'),
(uuid_generate_v4(), 'Apple', 'Medium apple', 182, 'g', 95, 0.5, 25, 0.3, 4.4, 'https://placehold.co/400x400?text=Apple', true, 'fruit'),
(uuid_generate_v4(), 'Strawberries', 'Fresh strawberries', 100, 'g', 32, 0.7, 7.7, 0.3, 2, 'https://placehold.co/400x400?text=Strawberries', true, 'fruit'),
(uuid_generate_v4(), 'Blueberries', 'Fresh blueberries', 100, 'g', 57, 0.7, 14.5, 0.3, 2.4, 'https://placehold.co/400x400?text=Blueberries', true, 'fruit');

-- Healthy Fats
INSERT INTO foods (id, name, description, serving_size, unit, calories, protein_g, carbs_g, fat_g, image_url, is_system, category) VALUES
(uuid_generate_v4(), 'Avocado', 'Fresh avocado', 100, 'g', 160, 2, 8.5, 14.7, 'https://placehold.co/400x400?text=Avocado', true, 'fat'),
(uuid_generate_v4(), 'Almonds', 'Raw almonds', 28, 'g', 164, 6, 6, 14, 'https://placehold.co/400x400?text=Almonds', true, 'fat'),
(uuid_generate_v4(), 'Olive Oil', 'Extra virgin olive oil', 14, 'ml', 119, 0, 0, 13.5, 'https://placehold.co/400x400?text=Olive+Oil', true, 'fat'),
(uuid_generate_v4(), 'Peanut Butter', 'Natural peanut butter', 32, 'g', 191, 7.7, 7, 16.4, 'https://placehold.co/400x400?text=Peanut+Butter', true, 'fat');

-- Dairy
INSERT INTO foods (id, name, description, serving_size, unit, calories, protein_g, carbs_g, fat_g, image_url, is_system, category) VALUES
(uuid_generate_v4(), 'Milk', 'Low-fat milk', 240, 'ml', 102, 8.2, 12.2, 2.4, 'https://placehold.co/400x400?text=Milk', true, 'dairy'),
(uuid_generate_v4(), 'Cheese', 'Cheddar cheese', 28, 'g', 113, 7, 0.9, 9.3, 'https://placehold.co/400x400?text=Cheese', true, 'dairy');

-- Beverages
INSERT INTO foods (id, name, description, serving_size, unit, calories, protein_g, carbs_g, fat_g, image_url, is_system, category) VALUES
(uuid_generate_v4(), 'Water', 'Plain water', 240, 'ml', 0, 0, 0, 0, 'https://placehold.co/400x400?text=Water', true, 'beverage'),
(uuid_generate_v4(), 'Green Tea', 'Unsweetened green tea', 240, 'ml', 2, 0.5, 0, 0, 'https://placehold.co/400x400?text=Green+Tea', true, 'beverage'),
(uuid_generate_v4(), 'Protein Shake', 'Whey protein shake', 240, 'ml', 120, 24, 3, 1.5, 'https://placehold.co/400x400?text=Protein+Shake', true, 'beverage');

-- Sample data: 30 common foods
SELECT 'Foods seeded successfully!' as message;
SELECT category, COUNT(*) as count FROM foods WHERE is_system = true GROUP BY category;
