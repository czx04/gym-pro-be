-- Sample Exercise Data
-- Run this manually after migrations to populate exercise library
-- Usage: psql -U gymadmin -d gym_pro_db -f migrations/seed_exercises.sql

-- Strength Training Exercises
INSERT INTO exercises (id, name, description, category, muscle_groups, equipment_needed, difficulty_level, calories_per_minute, is_active) VALUES
(uuid_generate_v4(), 'Push-ups', 'Classic bodyweight chest exercise', 'strength', '["chest", "triceps", "shoulders"]', '["none"]', 'beginner', 7.5, true),
(uuid_generate_v4(), 'Pull-ups', 'Upper body pulling exercise', 'strength', '["back", "biceps", "lats"]', '["pull-up bar"]', 'intermediate', 10.0, true),
(uuid_generate_v4(), 'Squats', 'Lower body compound exercise', 'strength', '["quads", "glutes", "hamstrings"]', '["none"]', 'beginner', 8.0, true),
(uuid_generate_v4(), 'Bench Press', 'Chest pressing exercise', 'strength', '["chest", "triceps", "shoulders"]', '["barbell", "bench"]', 'intermediate', 6.5, true),
(uuid_generate_v4(), 'Deadlift', 'Full body pulling exercise', 'strength', '["back", "glutes", "hamstrings", "core"]', '["barbell"]', 'advanced', 9.0, true),
(uuid_generate_v4(), 'Lunges', 'Single leg lower body exercise', 'strength', '["quads", "glutes", "hamstrings"]', '["none"]', 'beginner', 7.0, true),
(uuid_generate_v4(), 'Dumbbell Rows', 'Back rowing exercise', 'strength', '["back", "lats", "biceps"]', '["dumbbells"]', 'intermediate', 6.0, true),
(uuid_generate_v4(), 'Plank', 'Core stability exercise', 'strength', '["core", "abs", "shoulders"]', '["none"]', 'beginner', 5.0, true);

-- Cardio Exercises
INSERT INTO exercises (id, name, description, category, muscle_groups, equipment_needed, difficulty_level, calories_per_minute, is_active) VALUES
(uuid_generate_v4(), 'Running', 'Outdoor or treadmill running', 'cardio', '["legs", "cardiovascular"]', '["none"]', 'beginner', 11.5, true),
(uuid_generate_v4(), 'Cycling', 'Stationary or outdoor cycling', 'cardio', '["legs", "cardiovascular"]', '["bike"]', 'beginner', 9.0, true),
(uuid_generate_v4(), 'Jump Rope', 'Skipping rope cardio', 'cardio', '["legs", "shoulders", "cardiovascular"]', '["jump rope"]', 'intermediate', 13.0, true),
(uuid_generate_v4(), 'Burpees', 'Full body cardio exercise', 'cardio', '["full body", "cardiovascular"]', '["none"]', 'advanced', 12.5, true),
(uuid_generate_v4(), 'Mountain Climbers', 'Dynamic core and cardio', 'cardio', '["core", "legs", "cardiovascular"]', '["none"]', 'intermediate', 10.0, true);

-- Flexibility Exercises
INSERT INTO exercises (id, name, description, category, muscle_groups, equipment_needed, difficulty_level, calories_per_minute, is_active) VALUES
(uuid_generate_v4(), 'Hamstring Stretch', 'Seated or standing hamstring stretch', 'flexibility', '["hamstrings"]', '["none"]', 'beginner', 2.0, true),
(uuid_generate_v4(), 'Shoulder Stretch', 'Shoulder flexibility exercise', 'flexibility', '["shoulders"]', '["none"]', 'beginner', 2.0, true),
(uuid_generate_v4(), 'Yoga Flow', 'Dynamic yoga sequence', 'flexibility', '["full body"]', '["yoga mat"]', 'intermediate', 4.0, true),
(uuid_generate_v4(), 'Hip Flexor Stretch', 'Hip mobility exercise', 'flexibility', '["hips", "hip flexors"]', '["none"]', 'beginner', 2.5, true);

-- Stretching Exercises
INSERT INTO exercises (id, name, description, category, muscle_groups, equipment_needed, difficulty_level, calories_per_minute, is_active) VALUES
(uuid_generate_v4(), 'Cat-Cow Stretch', 'Spinal mobility exercise', 'stretching', '["back", "spine"]', '["none"]', 'beginner', 2.0, true),
(uuid_generate_v4(), 'Child Pose', 'Relaxing full body stretch', 'stretching', '["back", "shoulders", "hips"]', '["yoga mat"]', 'beginner', 1.5, true),
(uuid_generate_v4(), 'Quad Stretch', 'Standing quadriceps stretch', 'stretching', '["quads"]', '["none"]', 'beginner', 2.0, true);

-- Sample data: 20 exercises across all categories
SELECT 'Exercises seeded successfully!' as message;
SELECT category, COUNT(*) as count FROM exercises GROUP BY category;
