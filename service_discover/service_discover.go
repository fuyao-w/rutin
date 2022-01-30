package service_discover

import (
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
)

type ServiceDiscover interface {
	Register(addr string) error
	Unregister(addr string) error
	GetAddrSlice() (arrs []string)
	SetRegisterName(name string)
}

const (
	key = "service:discover:protocol:%s"
)

type RedisRegisterProtocol struct {
	name string
	rds  redigo.Conn
}

func NewRedisRegisterProtocol() (r *RedisRegisterProtocol, err error) {
	//return &RedisRegisterProtocol{
	//	rds: redis.NewClient(&redis.Options{
	//		Addr:            "127.0.0.1:6379",
	//		DB:              0,
	//		MaxRetries:      0,
	//		MinRetryBackoff: 0,
	//		MaxRetryBackoff: 0,
	//		DialTimeout:     0,
	//		ReadTimeout:     500 * time.Millisecond,
	//		WriteTimeout:    500 * time.Millisecond,
	//		PoolSize:        100,
	//		MinIdleConns:    10,
	//		MaxConnAge:      600,
	//		IdleTimeout:     600,
	//	}),
	//}, nil
	conn, err := redigo.Dial("tcp", "127.0.0.1:6379")
	return &RedisRegisterProtocol{
		rds: conn,
	}, err
}

func getKey(template string, args ...interface{}) string {
	return fmt.Sprintf(template, args...)
}
func (r *RedisRegisterProtocol) Register(addr string) error {
	_, err := r.rds.Do("sadd", getKey(key, r.name), addr)
	if err != nil {
		fmt.Println("redis register err ", err)
		return err
	}
	return nil
}

func (r *RedisRegisterProtocol) Unregister(addr string) error {
	_, err := r.rds.Do("srem", getKey(key, r.name), addr)
	if err != nil {
		fmt.Println("redis register err ", err)
		return err
	}
	return nil
}

func (r *RedisRegisterProtocol) GetAddrSlice() (arrs []string) {
	reply, err := redigo.Strings(r.rds.Do("SMembers", getKey(key, r.name)))
	if err != nil {
		fmt.Println("redis register err ", err)
	}
	return reply
}

func (r *RedisRegisterProtocol) SetRegisterName(name string) {
	r.name = name
}
