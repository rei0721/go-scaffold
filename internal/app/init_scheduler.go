package app

import (
	"fmt"
	"time"

	"github.com/rei0721/rei0721/pkg/scheduler"
)

func initScheduler(app *App) error {
	sched, err := scheduler.New(&scheduler.Config{
		PoolSize:       10000,
		ExpiryDuration: time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to create scheduler: %w", err)
	}
	app.Scheduler = sched
	app.Logger.Info("scheduler initialized", "poolSize", 10000)
	return nil
}
