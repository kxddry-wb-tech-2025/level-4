package handlers

import (
	"context"
	"log/slog"

	"github.com/kxddry/go-utils/pkg/logger/handlers/sl"
)

// HandleErrors logs errors coming from errCh
func HandleErrors(ctx context.Context, log *slog.Logger, errCh <-chan error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-errCh:
				if !ok {
					return
				}
				log.Error("found error", sl.Err(err))
			}
		}
	}()
}
