package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fuyao-w/rutin/consul"
	"github.com/fuyao-w/rutin/discovery"
	"github.com/fuyao-w/rutin/endpoint"
	"github.com/fuyao-w/rutin/iokit"
	"github.com/fuyao-w/rutin/rpc/internal/iosocket"
	"github.com/fuyao-w/rutin/rpc/internal/metadata"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go/token"
	"log"
	"reflect"
)

type (
	HandlerOption func(*HandlerOptions)

	HandlerOptions struct {
		HandlerName string
	}
	//Handler interface {
	//	Name() string
	//	Handler() interface{}
	//}
	Handler struct {
		name    string
		handler interface{}
	}
	Server interface {
		NewHandler(interface{}, ...HandlerOption) Handler

		// registe a Handler to server
		Handle(Handler) error

		// registe middlewares to a handler
		Use(...Plugin) Server

		Start() error
		Stop() error

		GetPaths() []string
	}

	methodType struct {
		method    reflect.Method
		ArgType   reflect.Type
		ReplyType reflect.Type
	}
	MethProc struct {
		meth   reflect.Method
		values []reflect.Value
	}
	Service struct {
		info *endpoint.ServiceInfo
		meth map[string]*methodType
		typ  reflect.Type
		rcvr reflect.Value
	}

	RpcServer struct {
		options Options
		//router
		serviceMap map[string]*Service
		plugins    []Plugin //熔断限流等插件
		server     *iosocket.Server
		shutdown   int32
		stop       chan struct{}
		once       sync.Once
		sync.RWMutex
	}

	ServerRegister interface {
		Name() string
	}
)

func getPath(service, meth string) string {
	return fmt.Sprintf("%s.%s", service, meth)
}

func WithServiceName(name string) HandlerOption {
	return func(options *HandlerOptions) {
		options.HandlerName = name
	}
}

func (r *RpcServer) NewHandler(handler interface{}, opts ...HandlerOption) Handler {
	var options HandlerOptions
	for _, opt := range opts {
		opt(&options)
	}
	return Handler{
		name: func() string {
			if options.HandlerName != "" {
				return options.HandlerName
			}
			return reflect.Indirect(reflect.ValueOf(handler)).Type().Name()
		}(),
		handler: handler,
	}
}

func (r *RpcServer) Handle(handler Handler) error {
	if len(handler.name) == 0 {
		return errors.New("handler don't have validate name")
	}

	_, ok := r.serviceMap[handler.name]
	if ok {
		return errors.New("service already exit")
	}
	typ := reflect.TypeOf(handler.handler)
	rcvr := reflect.ValueOf(handler.handler)

	if !isExportedOrBuiltinType(typ) {
		return errors.New("handler is not exported or builtin type")
	}
	service := Service{
		meth: suitableMethods(typ, false),
		typ:  typ,
		rcvr: rcvr,
		info: &endpoint.ServiceInfo{
			Name: handler.name,
			Addr: r.options.addr,
		},
	}
	r.Lock()
	defer r.Unlock()
	r.serviceMap[handler.name] = &service
	return nil
}

func (r *RpcServer) Use(plugin ...Plugin) Server {
	r.plugins = plugin
	return r
}

func (r *RpcServer) Start() error {
	var deferList []func()
	//1. 服务发现注册
	//2. 调用 iosocket.Server.Start()
	l, err := net.Listen("tcp4", r.options.addr.String())
	if err != nil {
		//logging.GenLogf("start rpc server on %s failed, %v", h.opts.Address, e)
		fmt.Printf("start rpc server on %s failed, %v\n", r.options.addr.String(), err)
		return err
	}
	reg := r.options.register
	for _, service := range r.serviceMap {
		if err = reg.Register(service.info); err != nil {
			log.Printf("service discover register failed :%s ,name :%s", err, r.options.ServiceName)
			l.Close()
			return err
		}
		deferList = append(deferList, func() {
			reg.Deregister(service.info)
		})
	}

	go func() {
		if err = r.server.Server.Start(l); err != nil {
			return
		}
	}()
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit
	for _, f := range deferList {
		f()
	}
	//注销服务发现
	//r.server.Server.Stop()
	log.Println("shut down")
	return nil
}

func (r *RpcServer) Stop() error {

	defer func() {
		close(r.stop)
	}()
	//注销服务发现
	r.server.Server.Stop()
	return nil
}

func (r *RpcServer) GetPaths() (paths []string) {
	for serviceName, service := range r.serviceMap {
		for methName := range service.meth {
			paths = append(paths, getPath(serviceName, methName))
		}
	}
	return
}

var (
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			}
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			}
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
			}
			continue
		}
		//Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			if reportErr {
				log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}

func (r *RpcServer) handleConnection(body []byte, wc io.WriteCloser) {
	var (
		desc = metadata.HandlerDesc{}
		pck  iokit.SeqPacket
	)
	//解析 协议包
	if err := pck.Decode(body); err != nil {
		fmt.Println("Decode err", err)
		return
	}
	//解析 rpc 请求元数据
	desc, err := metadata.Unmarshal(pck.Payload)
	if err != nil {
		log.Printf("handleConnection|Unmarshal err %s ,paylod :%s", err, string(body))
		return
	}
	handler, ok := r.serviceMap[desc.ServiceName]
	if !ok {
		log.Printf("handleConnection|service not found ,serviceName :%s ,methName :%s ,%s", desc.ServiceName, desc.MethName, string(pck.Payload))
		return
	}
	reply, err := r.call(*handler, desc)
	if err != nil {
		log.Printf("HandleConnection do err :%s ,%s", err, string(pck.Payload))
		return
	}
	//编码返回消息
	bytes, err := r.options.codec.Encode(reply.Interface())
	if err != nil {
		return
	}
	//iokit.SeqPacket{}
	desc.Response = bytes
	//编码 rpc 元数据
	pck.Payload, err = metadata.Marshal(&desc)
	if err != nil {
		return
	}
	//编码协议包
	bytes, err = pck.Encode()
	if err != nil {
		return
	}
	wc.Write(bytes)
}

func (h *RpcServer) serveRpc(remoteAddr string, request *iosocket.Context) (*iosocket.Context, error) {

	return nil, nil
}

func (r *RpcServer) call(handle Service, desc metadata.HandlerDesc) (reply reflect.Value, err error) {
	if meth, ok := handle.meth[desc.MethName]; !ok {
		return reply, errors.New("meth not found")
	} else {
		arg := reflect.New(meth.ArgType)
		reply = reflect.New(meth.ReplyType.Elem())
		if err := json.Unmarshal(desc.Param, arg.Interface()); err != nil {
			return reply, fmt.Errorf("do unmarshal err: %s", err.Error())
		}
		response := meth.method.Func.Call([]reflect.Value{handle.rcvr, arg.Elem(), reply})
		if r := response[0].Interface(); r != nil {
			err = r.(error)
		}

	}
	return reply, err
}

var DefaultRegisterCenter discovery.Register

func NewRpcServer(opts ...Option) Server {
	var (
		options = Options{
			register: consul.NewConsulDiscovery(),
		}
	)
	for _, opt := range opts {
		opt(&options)
	}
	s := &RpcServer{
		options:    options,
		serviceMap: map[string]*Service{},
		stop:       make(chan struct{}),
		once:       sync.Once{},
		RWMutex:    sync.RWMutex{},
	}
	s.server = iosocket.NewServer(s.handleConnection)

	return s
}
