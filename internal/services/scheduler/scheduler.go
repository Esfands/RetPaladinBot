package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
)

func Setup(ctx context.Context) (Service, error) {
	svc := &schedulerService{}

	svc.scheduler = gocron.NewScheduler(time.UTC)

	svc.scheduler.StartAsync()

	return svc, nil
}
