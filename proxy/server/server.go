package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fuyao-w/sd/net"
	"github.com/fuyao-w/sd/parse"
	"github.com/fuyao-w/sd/sd"
	"github.com/fuyao-w/sd/utils"
	"go/token"
	"log"
	netAddr "net"
	"reflect"
)

type (
	HandlerDesc struct {
		ServiceName string `json:"service_name"`
		MethName    string `json:"meth_name"`
		Param       []byte `json:"param"`
	}
	methodType struct {
		//sync.Mutex // protects counters
		method    reflect.Method
		ArgType   reflect.Type
		ReplyType reflect.Type
		//numCalls   uint
	}
	MethProc struct {
		meth   reflect.Method
		values []reflect.Value
	}
	Service struct {
		methRegister map[string]*methodType
		typ          reflect.Type
		rcvr         reflect.Value
	}
	Server struct {
		Name           string `json:"name"`
		Port           int    `json:"port"`
		addr           string `json:"addr"`
		sdComponent    sd.ServiceDiscover
		RegisterCenter RegisterCenter
	}
	RegisterCenter struct {
		Type string `json:"type"`
		Addr string `json:"addr"`
	}
	ServerRegister interface {
		Name() string
	}
)

var (
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
	serviceMap  = map[string]Service{}
)

func do(handle Service, desc HandlerDesc) (reply reflect.Value, err error) {
	if meth, ok := handle.methRegister[desc.MethName]; !ok {
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
func RegisterHandle(handle ServerRegister) {
	typ := reflect.TypeOf(handle)
	rcvr := reflect.ValueOf(handle)
	fmt.Println("RegisterHandle", handle.Name())
	serviceMap[handle.Name()] = Service{
		methRegister: suitableMethods(typ, false),
		typ:          typ,
		rcvr:         rcvr,
	}
}

func HandleConnection(conn netAddr.Conn) {
	var (
		parser = parse.ProtocolParser{}
		desc   = HandlerDesc{}
	)

	defer conn.Close()
	bytes, err := parser.Decode(conn) //解析tcp信息
	if err = json.Unmarshal(bytes, &desc); err != nil {
		log.Println("unmarshal err ", err)
		return
	}
	handler, ok := serviceMap[desc.ServiceName]
	if !ok {
		log.Printf("service not found :%s", desc.ServiceName)
		return
	}
	reply, err := do(handler, desc)
	if err != nil {
		log.Printf("HandleConnection do err :%s", err)
		return
	}

	bytes, err = parser.Encode(utils.GetJsonBytes(reply.Interface()))
	conn.Write(bytes)
}

func (s *Server) Server() {
	net.Server(s.addr, HandleConnection)
}

func (s *Server) Init() {
	var (
		sdComponent *sd.RedisRegisterProtocol
		err         error
	)
	s.addr = fmt.Sprintf("%s:%d", utils.GetIP(), s.Port)
	switch s.RegisterCenter.Type {
	case "redis":
		sdComponent, err = sd.NewRedisRegisterProtocol(s.RegisterCenter.Addr)
	case "consul":
		log.Fatal("consul register not implement")
		return
	}

	if err != nil {
		fmt.Println("init err ", err)
		return
	}
	sdComponent.SetRegisterName(s.Name)
	if err = sdComponent.Register(s.addr); err != nil {
		fmt.Println("consumer register ", err)
		return
	}
}
