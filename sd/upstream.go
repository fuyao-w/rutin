package sd

import (
	"context"
	"errors"
	"github.com/fuyao-w/rutin/core"
	"log"
)

func NewUpStream(fac PluginFactory, name string) core.Plugin {
	return core.Function(func(ctx context.Context, core core.Drive) {
		arrs := DefaultRegisterCenter.GetAddrSlice(name)
		if len(arrs) == 0 {
			log.Printf("client no upstream")
			core.AbortErr(errors.New("no upstream"))
			return
		}
		plugin, err := fac.Factory(arrs[0])
		if err != nil {
			log.Printf("NewUpStream|Factory err %s", err)
			core.AbortErr(err)
			return
		}
		plugin.Do(ctx, core)
		core.Next(ctx)
		if core.Err() != nil {
			log.Printf("NewUpStream|core.Err call failed %s", core.Err())
		}
	})
}
