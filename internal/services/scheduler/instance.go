package scheduler

import "github.com/go-co-op/gocron"

type Service interface {
	Scheduler() *gocron.Scheduler
}

type schedulerService struct {
	scheduler *gocron.Scheduler
}

func (s *schedulerService) Scheduler() *gocron.Scheduler {
	return s.scheduler
}
