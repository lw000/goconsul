package main

import (
	"demo/goconsul/protos"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

func main() {
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "greeter",
		}),
		// micro.Address("192.168.1.201"),
	)

	service.Init()

	greeter := Tpk.NewGreeterClient("greeter", service.Client())

	var count int64
	go func() {
		for {
			ctx := context.Background()
			resp, err := greeter.Ping(ctx, &Tpk.HelloRequest{Name: "Alice"})
			if err != nil {
				log.Println(err)
			}
			count = count + 1

			log.Println(count, resp)
		}
	}()

	select {}
}
