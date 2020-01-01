package ticketswitch

import (
	"context"
)

type key int

const (
	contextTrackingIdKey key = iota
)

// SetSessionTrackingID saves a tracking id into a context
func SetSessionTrackingID(ctx context.Context, trackingId string) context.Context {
	return context.WithValue(ctx, contextTrackingIdKey, trackingId)
}

// GetSessionTrackingID gets the tracking id from the context
func GetSessionTrackingID(ctx context.Context) (string, bool) {
	trackingId, ok := ctx.Value(contextTrackingIdKey).(string)
	return trackingId, ok
}
