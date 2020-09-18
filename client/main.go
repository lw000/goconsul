package main

import (
	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	lastIndex uint64
	g         sync.WaitGroup
)

func getService() {
	defer func() {
		if x := recover(); x != nil {

		}
		g.Done()
	}()
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500" // consul server

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Println("api new client is failed, err:", err)
		return
	}

	services, _ := client.Agent().Services()
	for _, value := range services {
		log.Println(value.Address, value.Port)
	}

	serves, metaInfo, err := client.Health().Service("serverNode", "v1", true, &consulapi.QueryOptions{
		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		log.Fatalf("error retrieving instances from Consul: %v", err)
	}
	lastIndex = metaInfo.LastIndex

	for _, serve := range serves {
		log.Println("service.Service.Address:", serve.Service.Address, "service.Service.Port:", serve.Service.Port)
	}
}

func checkHeath() {
	defer func() {
		if x := recover(); x != nil {

		}
		g.Done()
	}()
	cfg := consulapi.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	a, b, _ := client.Agent().AgentHealthServiceByID("serverNode")
	log.Println(a, b)
}

func consulKVTest() {
	defer func() {
		if x := recover(); x != nil {

		}
		g.Done()
	}()
	// 创建连接consul服务配置
	cfg := consulapi.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	// KV, put值
	values := "test"
	key := "go-consul-test/172.16.242.129:8100"
	client.KV().Put(&consulapi.KVPair{Key: key, Flags: 0, Value: []byte(values)}, nil)

	// KV get值
	data, _, _ := client.KV().Get(key, nil)
	log.Println(string(data.Value))

	// KV list
	datas, _, _ := client.KV().List("go", nil)
	for _, value := range datas {
		log.Println(value)
	}
	keys, _, _ := client.KV().Keys("go", "", nil)
	log.Println(keys)
}

func main() {

	g.Add(1)
	go checkHeath()
	g.Add(1)
	go getService()
	g.Add(1)
	go consulKVTest()

	g.Wait()
}
