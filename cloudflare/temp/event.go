package temp

import (
	"context"
	"errors"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
)

// Event represents information about the Cron that invoked this worker.
type Event struct {
	Cron string
	//ScheduledTime time.Time
}

func NewEvent(ctx context.Context) (*Event, error) {
	obj := runtimecontext.MustExtractRuntimeObj(ctx)
	if obj.IsUndefined() {
		return nil, errors.New("event is null")
	}

	//scheduledTimeVal := obj.Get("scheduledTime").Float()
	return &Event{
		Cron: "butt",
		//ScheduledTime: time.Unix(int64(scheduledTimeVal)/1000, 0).UTC(),
	}, nil
}
