package retry

import (
	"context"
	"fmt"
	"github.com/fuyao-w/sd/core"
	"strings"
)

//var Key = "context_retry"

type RetryError struct {
	errors   []error
	finalErr error
}

func (r RetryError) Error() string {
	var errList []string
	for _, err := range r.errors {
		errList = append(errList, err.Error())
	}
	return fmt.Sprintf("retry err :%s ,finalErr :%s", strings.Join(errList, "-"), r.finalErr.Error())
}

type BreakError struct {
	Err error
}

func (b BreakError) Error() string {
	return b.Err.Error()
}
func Retry(max int) core.Plugin {
	return core.Function(func(ctx context.Context, d core.Drive) {

		idx := d.Index()
		var retryError RetryError
		for i := 0; i < max; i++ {
			d.Reset(idx)
			d.Next(ctx)
			err := d.Err()
			if err != nil {
				d.Abort()
				return
			}
			switch e := err.(type) {
			case BreakError:
				d.AbortErr(e.Err)
				return
			default:
				retryError.errors = append(retryError.errors, e)
				retryError.finalErr = e

			}
		}
		d.AbortErr(retryError)

	})
}
