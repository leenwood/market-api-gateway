package rate

import (
	"sync"
	"time"
)

// limiter tracks sliding-window request counts per IP.
type limiter struct {
	mu        sync.Mutex
	counts    map[string]*windowCounter
	rps       int
	cleanupAt time.Time
}

type windowCounter struct {
	timestamps []time.Time
}

// RateLimiter implements port.RateLimiter with an in-memory sliding window.
type RateLimiter struct {
	l *limiter
}

func New(rps int) *RateLimiter {
	return &RateLimiter{
		l: &limiter{
			counts:    make(map[string]*windowCounter),
			rps:       rps,
			cleanupAt: time.Now().Add(time.Minute),
		},
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.l.mu.Lock()
	defer r.l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-time.Second)

	wc, ok := r.l.counts[ip]
	if !ok {
		wc = &windowCounter{}
		r.l.counts[ip] = wc
	}

	// Remove timestamps outside the 1-second window.
	valid := wc.timestamps[:0]
	for _, t := range wc.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	wc.timestamps = valid

	if len(wc.timestamps) >= r.l.rps {
		return false
	}
	wc.timestamps = append(wc.timestamps, now)

	// Periodic cleanup of idle IPs to prevent unbounded growth.
	if now.After(r.l.cleanupAt) {
		r.cleanup(cutoff)
		r.l.cleanupAt = now.Add(time.Minute)
	}

	return true
}

func (r *RateLimiter) cleanup(cutoff time.Time) {
	for ip, wc := range r.l.counts {
		if len(wc.timestamps) == 0 || wc.timestamps[len(wc.timestamps)-1].Before(cutoff) {
			delete(r.l.counts, ip)
		}
	}
}
