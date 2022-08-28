package client

import (
	"github.com/fuyao-w/circuit_breaker"
	"github.com/fuyao-w/rate_limit"
	"github.com/fuyao-w/rutin/discovery"
	"github.com/fuyao-w/rutin/rpc/codec"
	"log"
	"time"
)

type Options struct {
	//ServiceName string //简单通过服务发现名调用
	Codec          codec.RequestCodec
	PoolSize       int
	TTL            time.Duration
	RetryTimes     int
	Timeout        time.Duration
	rateLimit      rate_limit.RateLimit
	circuitBreaker *circuit_breaker.TwoStepCircuitBreaker
	dialer         dialer
	discovery      discovery.ServiceDiscover
}

type Option func(options *Options)

func newOptions(opts ...Option) (options Options) {
	options = Options{
		//Codec:      nil,
		PoolSize:   DefaultPoolSize,
		TTL:        DefaultPoolTTL,
		RetryTimes: 0,
		Timeout:    3 * time.Second,
		Codec:      &codec.JsonCodec{},
		rateLimit:  rate_limit.NewBucket(100000, 10000, time.Second),
		circuitBreaker: circuit_breaker.NewTwoStepCircuitBreaker(circuit_breaker.Options{
			Name:     "default",
			Interval: time.Minute,
			Timeout:  time.Second * 15,
			ReadyToTrip: func(counts circuit_breaker.Counts) bool {
				return float64(counts.ConsecutiveFailures/counts.TotalRequests) > 0.7

			},
			OnStateChange: func(name string, before, after circuit_breaker.State) {
				log.Println("circuit breaker state change", name, before, after)
			},
			IsSuccessful: func(err error) bool {
				return err == nil
			},
			Threshold: 100,
		}),
	}
	options.dialer = &defaultDialer{options: options}
	for _, opt := range opts {
		opt(&options)
	}
	return
}
