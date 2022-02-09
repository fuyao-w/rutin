package sd

import (
	"context"
	"github.com/fuyao-w/sd/core"
	"log"
)

func NewUpStream(fac PluginFactory, name string) core.Plugin {
	return core.Function(func(ctx context.Context, core core.Drive) {
		arrs := DefaultRegisterCenter.GetAddrSlice(name)
		if len(arrs) == 0 {
			log.Printf("client no upstream")
			return
		}
		plugin, err := fac.Factory(arrs[0])
		if err != nil {
			core.AbortErr(err)
			return
		}
		plugin.Do(ctx, core)
		core.Next(ctx)
		if core.Err() != nil {
			log.Printf("call failed", core.Err())
		}
	})
}
