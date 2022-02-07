package metadata

type HandlerDesc struct {
	ServiceName string `json:"service_name"`
	MethName    string `json:"meth_name"`
	Param       []byte `json:"param"`
}
