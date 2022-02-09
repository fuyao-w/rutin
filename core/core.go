package core

import "context"

type Drive interface {
	Use(...Plugin) Drive
	Next(context.Context)
	AbortErr(error)
	Abort()
	IsAborted() bool
	Err() error
	Copy() Drive
	Index() int
	Reset(idx int)
}

func New(plugins []Plugin) Drive {
	return &Driver{
		index:   -1,
		plugins: plugins,
	}
}

type Driver struct {
	plugins []Plugin
	index   int
	err     error
}

// Deprecated: Copy, just use Index()
func (d *Driver) Copy() Drive {
	dup := &Driver{}
	dup.index = d.index
	dup.plugins = append(d.plugins[:0:0], d.plugins...)
	return dup
}

func (d *Driver) Use(ps ...Plugin) Drive {
	d.plugins = append(d.plugins, ps...)
	return d
}

func (d *Driver) Next(ctx context.Context) {
	d.index++
	for s := len(d.plugins); d.index < s; d.index++ {
		d.plugins[d.index].Do(ctx, d)
	}
}

func (d *Driver) Abort() {
	d.index = len(d.plugins)
}

func (d *Driver) AbortErr(err error) {
	d.Abort()
	d.err = err
}

func (d *Driver) Err() error {
	return d.err
}

func (d *Driver) IsAborted() bool {
	return d.index >= len(d.plugins)
}

func (d *Driver) Index() int {
	return d.index
}

func (d *Driver) Reset(n int) {
	d.index = n
	d.err = nil
}

type Plugin interface {
	Do(ctx context.Context, Driver Drive)
}

type Function func(ctx context.Context, Driver Drive)

func (f Function) Do(ctx context.Context, d Drive) { f(ctx, d) }
