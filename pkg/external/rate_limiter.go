package external

import (
	"fmt"
	"sync"
	"time"
)

type DailyLimiter struct {
	mu      sync.Mutex
	count   int
	resetAt time.Time
	limit   int
}

func NewDailyLimiter(limit int) *DailyLimiter {
	return &DailyLimiter{
		limit:   limit,
		resetAt: getTomorrowMidnight(),
	}
}

func getTomorrowMidnight() time.Time {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location())
}

func (l *DailyLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.After(l.resetAt) {
		l.count = 0
		l.resetAt = getTomorrowMidnight()
	}

	if l.count >= l.limit {
		return false
	}

	l.count++
	fmt.Println("OMDb counter:", l.count)
	return true
}
