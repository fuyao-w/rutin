package upstream

import (
	"context"
	"errors"
	"github.com/fuyao-w/rutin/core"
	"github.com/fuyao-w/rutin/discovery"
	"github.com/fuyao-w/rutin/load_balance"
	"log"
)

func NewUpStream(fac discovery.PluginFactory, getCollection func() discovery.Collection) core.Plugin {
	return core.Function(func(ctx context.Context, core core.Drive) {
		instance := load_balance.GetLocalLoadBalancer().GetPicker(getCollection()).Pick()
		if instance == nil {
			log.Printf("client no upstream")
			core.AbortErr(errors.New("no upstream"))
			return
		}
		plugin, err := fac.Factory(instance.GetAddr().String())
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
