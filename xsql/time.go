package xsql

import (
	"context"
	"time"
)

func Shrink(ctx context.Context, duration time.Duration) (time.Duration, context.Context, context.CancelFunc) {
	if duration == 0 {
		return 0, ctx, func() {}
	}
	if deadline, ok := ctx.Deadline(); ok {
		if left := time.Until(deadline); left < duration {
			return left, ctx, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(ctx, duration)
	return duration, ctx, cancel
}
