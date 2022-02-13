package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fuyao-w/rutin/iokit"
	"github.com/fuyao-w/rutin/rpc/internal/iosocket"
	"github.com/fuyao-w/rutin/rpc/internal/metadata"
	"github.com/fuyao-w/rutin/sd"
	"io"
	"net"
	"sync"

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
		name string
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
		name: handler.name,
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
	//1. 服务发现注册
	//2. 调用 iosocket.Server.Start()
	l, err := net.Listen("tcp4", r.options.Address)
	if err != nil {
		//logging.GenLogf("start rpc server on %s failed, %v", h.opts.Address, e)
		fmt.Printf("start rpc server on %s failed, %v\n", r.options.Address, err)
		return err
	}
	for service := range r.serviceMap {
		if err = sd.DefaultRegisterCenter.Register(service, r.options.Address); err != nil {
			log.Printf("service discover register failed :%s ,name :%s", err, r.options.ServiceName)
			l.Close()
			return err
		}
	}

	if err = r.server.Server.Start(l); err != nil {
		return err
	}
	<-r.stop
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

func NewRpcServer(opts ...Option) Server {
	var options Options
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
