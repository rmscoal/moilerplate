package rater

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	_defaultRateForEachClient rate.Limit = 2
	_defaultBurstByEachClient int        = 2

	_defaultIntervalEvaluation     = time.Minute
	_defaultClientDeletionInterval = 3 * time.Minute
)

var (
	raterSingleInstance *Rater
	once                sync.Once
)

type Rater struct {
	clients map[string]*client
	mu      sync.Mutex

	rateForEachClient  rate.Limit
	burstForEachClient int

	intervalForEvaluation  time.Duration
	clientDeletionInterval time.Duration
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func GetRater(ctx context.Context, opts ...Option) *Rater {
	if raterSingleInstance == nil {
		once.Do(func() {
			raterSingleInstance = &Rater{
				clients: make(map[string]*client),

				rateForEachClient:  _defaultRateForEachClient,
				burstForEachClient: _defaultBurstByEachClient,

				intervalForEvaluation:  _defaultIntervalEvaluation,
				clientDeletionInterval: _defaultClientDeletionInterval,
			}

			for _, opt := range opts {
				opt(raterSingleInstance)
			}

			raterSingleInstance.evaluateInBackground(ctx)
		})
	}

	return raterSingleInstance
}

func (rater *Rater) evaluateInBackground(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(rater.intervalForEvaluation)

				rater.mu.Lock()
				for ip, client := range rater.clients {
					if time.Since(client.lastSeen) > rater.clientDeletionInterval {
						delete(rater.clients, ip)
					}
				}
				rater.mu.Unlock()
			}
		}
	}()
}

func (rater *Rater) GetClient(ip string) (*client, bool) {
	rater.mu.Lock()
	defer rater.mu.Unlock()

	client, found := rater.clients[ip]
	return client, found
}

func (rater *Rater) AddNewClient(ip string) *client {
	rater.mu.Lock()
	defer rater.mu.Unlock()

	client := &client{limiter: rate.NewLimiter(rater.rateForEachClient, rater.burstForEachClient), lastSeen: time.Now()}
	rater.clients[ip] = client
	return client
}

func (rater *Rater) IsClientAllowed(client *client) bool {
	rater.mu.Lock()
	defer rater.mu.Unlock()

	client.lastSeen = time.Now()
	return client.limiter.Allow()
}
