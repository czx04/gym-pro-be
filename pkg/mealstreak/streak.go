package mealstreak

import (
	"sort"
	"time"
)

// VNLoc is Asia/Ho_Chi_Minh (no DST).
var VNLoc = time.FixedZone("ICT", 7*3600)

// TodayYYYYMMDD returns today's calendar date in Vietnam as YYYY-MM-DD.
func TodayYYYYMMDD(now time.Time) string {
	t := now.In(VNLoc)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, VNLoc).Format("2006-01-02")
}

// YesterdayYYYYMMDD returns yesterday's date in Vietnam.
func YesterdayYYYYMMDD(now time.Time) string {
	t := now.In(VNLoc).AddDate(0, 0, -1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, VNLoc).Format("2006-01-02")
}

// AddDays shifts a YYYY-MM-DD string by delta days in Vietnam calendar.
func AddDays(ymd string, delta int) string {
	t, err := time.ParseInLocation("2006-01-02", ymd, VNLoc)
	if err != nil {
		return ymd
	}
	t = t.AddDate(0, 0, delta)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, VNLoc).Format("2006-01-02")
}

// CurrentStreak counts consecutive logged days ending at today or yesterday (VN).
func CurrentStreak(sortedUniqueDates []string, now time.Time) int {
	if len(sortedUniqueDates) == 0 {
		return 0
	}
	set := make(map[string]struct{}, len(sortedUniqueDates))
	for _, d := range sortedUniqueDates {
		set[d] = struct{}{}
	}
	today := TodayYYYYMMDD(now)
	yesterday := YesterdayYYYYMMDD(now)

	var start string
	if _, ok := set[today]; ok {
		start = today
	} else if _, ok := set[yesterday]; ok {
		start = yesterday
	} else {
		return 0
	}

	n := 0
	d := start
	for {
		if _, ok := set[d]; !ok {
			break
		}
		n++
		d = AddDays(d, -1)
	}
	return n
}

// LongestStreak returns the maximum length of any consecutive run in history.
func LongestStreak(sortedUniqueDates []string) int {
	if len(sortedUniqueDates) == 0 {
		return 0
	}
	dates := append([]string(nil), sortedUniqueDates...)
	sort.Strings(dates)
	best := 1
	run := 1
	for i := 1; i < len(dates); i++ {
		if AddDays(dates[i-1], 1) == dates[i] {
			run++
			if run > best {
				best = run
			}
		} else {
			run = 1
		}
	}
	return best
}
