package utils

import (
	"context"
	"example_project/internal/leveledlog"
	"sync"
	"time"
)

var ctxLogger = leveledlog.Logger.WithID("unprocessed_event_worker")

func HandleUnprocessedEvent(ctx context.Context, wg *sync.WaitGroup, retryPeriod time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(retryPeriod)
	ctxLogger.Info("unprocessed event worker started successfully")
	for {
		select {
		case <-ctx.Done():
			ctxLogger.Info("context cancelled, exiting unprocessed event worker")
			return
		case <-ticker.C:
			ctxLogger.Info("retrieving unprocessed transactions")

			wg.Add(1)
			HandleEvent(wg)

		}
	}
}
func HandleEvent(wg *sync.WaitGroup) {
	defer wg.Done()
}
