package consul

import (
	"fmt"
	"github.com/fuyao-w/rutin/discovery"
	"github.com/fuyao-w/rutin/endpoint"
	"github.com/hashicorp/consul/api"
	"log"
	"net"
	"sync"
	"time"
)

var (
	client *api.Client
	once   sync.Once
)

type DiscoveryConsul struct{}

func NewConsulDiscovery() *DiscoveryConsul {
	initClient()
	return &DiscoveryConsul{}
}

func (c *DiscoveryConsul) Register(info *endpoint.ServiceInfo) (err error) {
	id := buildDiscoveryKey(info)
	addr := info.Addr.(*net.TCPAddr)
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      id,
		Name:    info.Name,
		Tags:    info.Tags,
		Port:    addr.Port,
		Address: addr.IP.String(),
		//SocketPath:        "",

		EnableTagOverride: false,
		Meta: map[string]string{
			"register_at": fmt.Sprintf("%d", time.Now().UnixMilli()),
		},
		//Weights:   nil,

		Check: &api.AgentServiceCheck{
			Name:              "check" + id,
			DockerContainerID: "",
			Interval:          "5s",
			Timeout:           "3s",
			TCP:               addr.String(),
			Notes:             "rutin register",
			Status:            api.HealthPassing,
			//SuccessBeforePassing:   0,
			//FailuresBeforeWarning:  1,
			//FailuresBeforeCritical: 2,
			//DeregisterCriticalServiceAfter: "", // 默认半分钟
		},
		Namespace: "",
	})

	if err != nil {
		panic(fmt.Sprintf("consul register err :%s", err))
	}
	return
}
func buildDiscoveryKey(info *endpoint.ServiceInfo) string {
	return info.Addr.String()
}
func (c *DiscoveryConsul) Deregister(info *endpoint.ServiceInfo) (err error) {
	id := buildDiscoveryKey(info)
	log.Println("Unregister", id)
	if err = client.Agent().ServiceDeregister(id); err != nil {
		log.Printf("consul deregister err :%s", err)
	}
	return
}

func (c *DiscoveryConsul) GetCollection(info *endpoint.ServiceInfo) (collection discovery.Collection) {
	var list []discovery.Instance
	collection = new(insCollection)
	entryList, _, err := client.Health().Service(info.Name, "", true, nil)
	if err != nil {
		log.Printf("GetAddrSlice err :%s", err)
		return collection
	}
	for _, entry := range entryList {
		s := entry.Service
		list = append(list, &ins{
			addr: &net.TCPAddr{
				IP:   net.ParseIP(s.Address),
				Port: s.Port,
			},
			tags:   s.Tags,
			weight: s.Weights.Passing,
		})
	}
	return &insCollection{
		list: list,
	}
}

func initClient() {
	once.Do(func() {
		var err error
		client, err = api.NewClient(api.DefaultConfig())
		if err != nil {
			panic(fmt.Errorf("consul init client err :%s", err))
		}
	})
}
