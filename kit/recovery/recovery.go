package recovery

import (
	"context"
	"github.com/fuyao-w/sd/core"
	"log"
	"runtime/debug"
)

func Recover() core.Plugin {
	return core.Function(func(ctx context.Context, core core.Drive) {
		defer func() {
			if p := recover(); p != nil {
				debug.PrintStack()
				log.Printf("recover panic !!! ,%v", p)
			}
		}()
		core.Next(ctx)
	})
}
