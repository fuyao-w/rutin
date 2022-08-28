package circurt_breaker

import (
	"context"
	"github.com/fuyao-w/circuit_breaker"
	"github.com/fuyao-w/rutin/core"
)

func CircuitBreaker(breaker *circuit_breaker.TwoStepCircuitBreaker) core.Plugin {
	return core.Function(func(ctx context.Context, core core.Drive) {
		success, err := breaker.IsAllow()
		if err == circuit_breaker.ErrCircuitBreaker {
			core.AbortErr(err)
			return
		}
		core.Next(ctx)
		success(core.Err() != nil)
	})
}
