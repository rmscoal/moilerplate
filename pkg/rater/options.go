package rater

import (
	"math"
	"time"

	"golang.org/x/time/rate"
)

// Option -.
type Option func(*Rater)

func RegisterRateLimitForEachClient(lim int) Option {
	return func(r *Rater) {
		if lim >= 0 && float64(lim) <= math.MaxFloat64 {
			r.rateForEachClient = rate.Limit(lim)
		}
	}
}

func RegisterBurstLimitForEachClient(burst int) Option {
	return func(r *Rater) {
		if burst >= 0 && float64(burst) <= math.MaxFloat64 {
			r.burstForEachClient = burst
		}
	}
}

func RegisterEvaluationInterval(t time.Duration) Option {
	return func(r *Rater) {
		if t > time.Duration(0) {
			r.intervalForEvaluation = t
		}
	}
}

func RegisterDeletionTime(t time.Duration) Option {
	return func(r *Rater) {
		if t > time.Duration(0) {
			r.clientDeletionInterval = t
		}
	}
}
