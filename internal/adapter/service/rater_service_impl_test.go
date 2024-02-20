package service

import (
	"context"
	"testing"
	"time"

	"github.com/rmscoal/moilerplate/pkg/rater"
	"github.com/stretchr/testify/assert"
)

func TestRaterService_IsClientAllowed_One_Burst(t *testing.T) {
	rt := rater.GetRater(context.Background(),
		rater.RegisterRateLimitForEachClient(2),
		rater.RegisterBurstLimitForEachClient(2),
		rater.RegisterEvaluationInterval(time.Minute),
		rater.RegisterDeletionTime(time.Minute),
	)

	service := NewRaterService(rt)

	assert.True(t, service.IsClientAllowed(context.Background(), "some_ip"))
	assert.True(t, service.IsClientAllowed(context.Background(), "some_ip"))
	assert.False(t, service.IsClientAllowed(context.Background(), "some_ip"))
}
