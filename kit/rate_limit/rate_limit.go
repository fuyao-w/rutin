package rate_limit

import (
	"context"
	"github.com/fuyao-w/rate_limit"
	"github.com/fuyao-w/rutin/core"
)

func RateLimit(limiter rate_limit.RateLimit) core.Plugin {
	return core.Function(func(ctx context.Context, core core.Drive) {
		_ = limiter.Take(1)
		core.Next(ctx)
	})
}
