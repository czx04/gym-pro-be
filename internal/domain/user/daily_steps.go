package user

import (
	"time"

	"github.com/google/uuid"
)

const DailyStepsSourceAppleHealth = "apple_health"

// DailySteps is one per user per local calendar day.
type DailySteps struct {
	UserID    uuid.UUID `json:"user_id"`
	Date      time.Time `json:"date"`
	Steps     int       `json:"steps"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DailyStepsPoint struct {
	Date  string `json:"date"` // YYYY-MM-DD
	Steps int    `json:"steps"`
}

