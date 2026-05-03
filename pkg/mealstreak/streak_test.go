package mealstreak

import (
	"testing"
	"time"
)

func TestCurrentStreak(t *testing.T) {
	loc := VNLoc
	mustParse := func(s string) time.Time {
		tm, err := time.ParseInLocation("2006-01-02 15:04:05", s+" 12:00:00", loc)
		if err != nil {
			t.Fatal(err)
		}
		return tm
	}
	dates := []string{"2026-03-27", "2026-03-28", "2026-03-29"}
	now := mustParse("2026-03-29")
	if got := CurrentStreak(dates, now); got != 3 {
		t.Fatalf("got %d want 3", got)
	}
	// no log today or yesterday
	if got := CurrentStreak([]string{"2026-03-25"}, mustParse("2026-03-29")); got != 0 {
		t.Fatalf("got %d want 0", got)
	}
	// only yesterday counts
	if got := CurrentStreak([]string{"2026-03-28"}, mustParse("2026-03-29")); got != 1 {
		t.Fatalf("got %d want 1", got)
	}
}

func TestLongestStreak(t *testing.T) {
	if got := LongestStreak([]string{"2026-01-01", "2026-01-02", "2026-01-05", "2026-01-06"}); got != 2 {
		t.Fatalf("got %d want 2", got)
	}
	if got := LongestStreak([]string{}); got != 0 {
		t.Fatalf("got %d want 0", got)
	}
}
