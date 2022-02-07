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

type Driver struct {
	plugins []Plugin
	idx     int
}

func (d Driver) Use(plugin ...Plugin) Drive {
	panic("implement me")
}

func (d Driver) Next(ctx context.Context) {
	panic("implement me")
}

func (d Driver) AbortErr(err error) {
	panic("implement me")
}

func (d Driver) Abort() {
	panic("implement me")
}

func (d Driver) IsAborted() bool {
	panic("implement me")
}

func (d Driver) Err() error {
	panic("implement me")
}

func (d Driver) Copy() Drive {
	panic("implement me")
}

func (d Driver) Index() int {
	panic("implement me")
}

func (d Driver) Reset(idx int) {
	panic("implement me")
}

func New(plugins []Plugin) Driver {
	return Driver{
		plugins: plugins,
		idx:     -1,
	}
}

type Plugin interface {
	Do(ctx context.Context, core Drive)
}

type Function func(ctx context.Context, core Drive)

func (f Function) Do(context.Context, Drive) {}
