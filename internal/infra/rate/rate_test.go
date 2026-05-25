package rate_test

import (
	"testing"
	"time"

	"market-api-gateway/internal/infra/rate"
)

func TestRateLimiter_Allow_UnderLimit(t *testing.T) {
	rl := rate.New(5)
	for i := range 5 {
		if !rl.Allow("1.2.3.4") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
}

func TestRateLimiter_Allow_ExceedLimit(t *testing.T) {
	rl := rate.New(3)
	for range 3 {
		rl.Allow("1.2.3.4")
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("4th request should be denied")
	}
}

func TestRateLimiter_Allow_SeparateIPs(t *testing.T) {
	rl := rate.New(2)
	rl.Allow("1.1.1.1")
	rl.Allow("1.1.1.1")

	// Different IP should still have its own quota.
	if !rl.Allow("2.2.2.2") {
		t.Fatal("different IP should be allowed")
	}
}

func TestRateLimiter_Allow_WindowExpiry(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timing-sensitive test in short mode")
	}
	rl := rate.New(2)
	rl.Allow("1.2.3.4")
	rl.Allow("1.2.3.4")

	// Deny at limit.
	if rl.Allow("1.2.3.4") {
		t.Fatal("should be denied at limit")
	}

	// After the 1-second sliding window, quota resets.
	time.Sleep(1100 * time.Millisecond)

	if !rl.Allow("1.2.3.4") {
		t.Fatal("should be allowed after window expiry")
	}
}
