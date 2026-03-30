package bootstrap

import (
	"context"
	"time"

	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/domain/meal"
	"gym-pro-2026-ptit/internal/infrastructure/expo"
	"gym-pro-2026-ptit/internal/infrastructure/logger"
	"gym-pro-2026-ptit/pkg/mealstreak"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

const mealReminderTitle = "Nhắc log meal"
const mealReminderBody = "Đừng quên ghi lại bữa ăn hôm nay để giữ streak."

// RegisterMealReminderCron schedules daily 17:00 Asia/Ho_Chi_Minh push for users without a meal log today.
func RegisterMealReminderCron(lc fx.Lifecycle, pushRepo meal.PushTokenRepository, cfg *config.Config) {
	if cfg.Expo.AccessToken == "" {
		logger.Info("meal reminder cron disabled", "reason", "EXPO_ACCESS_TOKEN is empty")
		return
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		logger.Error("failed to load Asia/Ho_Chi_Minh", "err", err)
		return
	}

	c := cron.New(cron.WithLocation(loc))
	_, err = c.AddFunc("0 17 * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		today := time.Now().In(mealstreak.VNLoc)
		todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, mealstreak.VNLoc)

		rows, err := pushRepo.ListTokensUsersWithoutMealOnDate(ctx, todayDate)
		if err != nil {
			logger.Error("meal reminder: list tokens", "err", err)
			return
		}
		if len(rows) == 0 {
			return
		}

		const batchSize = 99
		for i := 0; i < len(rows); i += batchSize {
			end := i + batchSize
			if end > len(rows) {
				end = len(rows)
			}
			var msgs []expo.Message
			for _, row := range rows[i:end] {
				msgs = append(msgs, expo.Message{
					To:    row.ExpoPushToken,
					Title: mealReminderTitle,
					Body:  mealReminderBody,
					Sound: "default",
				})
			}
			if err := expo.SendMessages(ctx, cfg.Expo.AccessToken, msgs); err != nil {
				logger.Error("meal reminder: expo send", "err", err, "batch_start", i)
			}
		}
		logger.Info("meal reminder sent", "recipients", len(rows))
	})
	if err != nil {
		logger.Error("meal reminder cron add func", "err", err)
		return
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			c.Start()
			logger.Info("meal reminder cron started", "schedule", "17:00 Asia/Ho_Chi_Minh daily")
			return nil
		},
		OnStop: func(_ context.Context) error {
			ctx := c.Stop()
			<-ctx.Done()
			return nil
		},
	})
}
