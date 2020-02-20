package main

import (
	"context"
	Tpk "demo/goconsul/protos"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"time"
)

type Greeter struct {
}

func (h *Greeter) Ping(ctx context.Context, req *Tpk.HelloRequest, rep *Tpk.HelloResponse) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		log.Printf("%+v", req)
		rep.Greeting = "Hello " + req.Name
	}
	return nil
}

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "greeter",
		}),
		micro.Registry(reg),
		micro.RegisterInterval(time.Second*time.Duration(20)),
		micro.RegisterTTL(time.Second*time.Duration(30)),
		// micro.Address("192.168.1.201"),
	)

	service.Init()

	Tpk.RegisterGreeterHandler(service.Server(), &Greeter{})

	if er := service.Run(); er != nil {
		log.Panic(er)
	}
}
