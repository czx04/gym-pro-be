package postgres

import (
	"context"
	"gym-pro-2026-ptit/internal/domain/user"
	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, role, is_active,
			oauth_provider, oauth_id,
			name, bio, avatar_url, date_of_birth, gender,
			height_cm, weight_kg, fitness_goal, activity_level,
			daily_calorie_target, protein_target_g, carbs_target_g, fat_target_g,
			privacy_settings, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		)
	`

	_, err := r.db.Exec(ctx, query,
		u.ID, u.Email, u.PasswordHash, u.Role, u.IsActive,
		u.OAuthProvider, u.OAuthID,
		u.Name, u.Bio, u.AvatarURL, u.DateOfBirth, u.Gender,
		u.HeightCm, u.WeightKg, u.FitnessGoal, u.ActivityLevel,
		u.DailyCalorieTarget, u.ProteinTargetG, u.CarbsTargetG, u.FatTargetG,
		u.PrivacySettings, u.CreatedAt, u.UpdatedAt,
	)

	if err != nil {
		return errors.DatabaseError("create user", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active,
			   oauth_provider, oauth_id,
			   name, bio, avatar_url, date_of_birth, gender,
			   height_cm, weight_kg, fitness_goal, activity_level,
			   daily_calorie_target, protein_target_g, carbs_target_g, fat_target_g,
			   privacy_settings, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u user.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive,
		&u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.Bio, &u.AvatarURL, &u.DateOfBirth, &u.Gender,
		&u.HeightCm, &u.WeightKg, &u.FitnessGoal, &u.ActivityLevel,
		&u.DailyCalorieTarget, &u.ProteinTargetG, &u.CarbsTargetG, &u.FatTargetG,
		&u.PrivacySettings, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("user")
		}
		return nil, errors.DatabaseError("get user by id", err)
	}

	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active,
			   oauth_provider, oauth_id,
			   name, bio, avatar_url, date_of_birth, gender,
			   height_cm, weight_kg, fitness_goal, activity_level,
			   daily_calorie_target, protein_target_g, carbs_target_g, fat_target_g,
			   privacy_settings, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u user.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive,
		&u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.Bio, &u.AvatarURL, &u.DateOfBirth, &u.Gender,
		&u.HeightCm, &u.WeightKg, &u.FitnessGoal, &u.ActivityLevel,
		&u.DailyCalorieTarget, &u.ProteinTargetG, &u.CarbsTargetG, &u.FatTargetG,
		&u.PrivacySettings, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("user")
		}
		return nil, errors.DatabaseError("get user by email", err)
	}

	return &u, nil
}

func (r *userRepository) GetByOAuth(ctx context.Context, provider, oauthID string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active,
			   oauth_provider, oauth_id,
			   name, bio, avatar_url, date_of_birth, gender,
			   height_cm, weight_kg, fitness_goal, activity_level,
			   daily_calorie_target, protein_target_g, carbs_target_g, fat_target_g,
			   privacy_settings, created_at, updated_at
		FROM users
		WHERE oauth_provider = $1 AND oauth_id = $2
	`

	var u user.User
	err := r.db.QueryRow(ctx, query, provider, oauthID).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive,
		&u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.Bio, &u.AvatarURL, &u.DateOfBirth, &u.Gender,
		&u.HeightCm, &u.WeightKg, &u.FitnessGoal, &u.ActivityLevel,
		&u.DailyCalorieTarget, &u.ProteinTargetG, &u.CarbsTargetG, &u.FatTargetG,
		&u.PrivacySettings, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("user")
		}
		return nil, errors.DatabaseError("get user by oauth", err)
	}

	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users SET
			email = $2, name = $3, bio = $4, avatar_url = $5,
			date_of_birth = $6, gender = $7, height_cm = $8, weight_kg = $9,
			fitness_goal = $10, activity_level = $11,
			daily_calorie_target = $12, protein_target_g = $13,
			carbs_target_g = $14, fat_target_g = $15,
			privacy_settings = $16, updated_at = $17
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		u.ID, u.Email, u.Name, u.Bio, u.AvatarURL,
		u.DateOfBirth, u.Gender, u.HeightCm, u.WeightKg,
		u.FitnessGoal, u.ActivityLevel,
		u.DailyCalorieTarget, u.ProteinTargetG, u.CarbsTargetG, u.FatTargetG,
		u.PrivacySettings, u.UpdatedAt,
	)

	if err != nil {
		return errors.DatabaseError("update user", err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return errors.DatabaseError("delete user", err)
	}

	if result.RowsAffected() == 0 {
		return errors.NotFound("user")
	}

	return nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, id uuid.UUID, input user.UpdateProfileInput) error {
	query := `
		UPDATE users SET
			name = COALESCE($2, name),
			bio = COALESCE($3, bio),
			avatar_url = COALESCE($4, avatar_url),
			date_of_birth = COALESCE($5, date_of_birth),
			gender = COALESCE($6, gender),
			height_cm = COALESCE($7, height_cm),
			weight_kg = COALESCE($8, weight_kg),
			fitness_goal = COALESCE($9, fitness_goal),
			activity_level = COALESCE($10, activity_level),
			daily_calorie_target = COALESCE($11, daily_calorie_target),
			protein_target_g = COALESCE($12, protein_target_g),
			carbs_target_g = COALESCE($13, carbs_target_g),
			fat_target_g = COALESCE($14, fat_target_g),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		id,
		input.Name,
		input.Bio,
		input.AvatarURL,
		input.DateOfBirth,
		input.Gender,
		input.HeightCm,
		input.WeightKg,
		input.FitnessGoal,
		input.ActivityLevel,
		input.DailyCalorieTarget,
		input.ProteinTargetG,
		input.CarbsTargetG,
		input.FatTargetG,
	)

	if err != nil {
		return errors.DatabaseError("update profile", err)
	}

	if result.RowsAffected() == 0 {
		return errors.NotFound("user")
	}

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, passwordHash)
	if err != nil {
		return errors.DatabaseError("update password", err)
	}

	if result.RowsAffected() == 0 {
		return errors.NotFound("user")
	}

	return nil
}

func (r *userRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, errors.DatabaseError("check user exists", err)
	}

	return exists, nil
}
