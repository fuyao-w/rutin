package sd

import (
	"fmt"
	"github.com/fuyao-w/rutin/core"
	redigo "github.com/garyburd/redigo/redis"
	"log"
)

type ServiceDiscover interface {
	Register(name, addr string) error
	Unregister(name, addr string) error
	GetAddrSlice(name string) (arrs []string)
}

type PluginFactory interface {
	Factory(host string) (core.Plugin, error)
}

const (
	key = "service:discover:protocol:%s"
)

type RedisRegisterProtocol struct {
	rds  redigo.Conn
	addr string
}

func NewRedisRegisterProtocol(addr string) (r *RedisRegisterProtocol, err error) {
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

	conn, err := redigo.Dial("tcp", addr)
	return &RedisRegisterProtocol{
		rds:  conn,
		addr: addr,
	}, err
}

func getKey(template string, args ...interface{}) string {
	return fmt.Sprintf(template, args...)
}
func (r *RedisRegisterProtocol) Register(name, addr string) error {
	_, err := r.rds.Do("sadd", getKey(key, name), addr)
	if err != nil {
		fmt.Println("redis register err ", err)
		return err
	}
	return nil
}

func (r *RedisRegisterProtocol) Unregister(name, addr string) error {
	_, err := r.rds.Do("srem", getKey(key, name), addr)
	if err != nil {
		fmt.Println("redis register err ", err)
		return err
	}
	return nil
}

func (r *RedisRegisterProtocol) GetAddrSlice(name string) (arrs []string) {
	return []string{"127.0.0.1:10000"}
	conn, err := redigo.Dial("tcp", r.addr)
	if err != nil {
		log.Printf("GetAddrSlice|Dial %s", err)
		return nil
	}
	defer func() {
		conn.Close()
	}()
	reply, err := redigo.Strings(conn.Do("SMEMBERS", getKey(key, name)))
	if err != nil {
		log.Printf("redis register err: %s", err)
	}
	return reply
}
