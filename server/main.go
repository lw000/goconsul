package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	count     int64 = 1
	checkPort       = 8080
)

func localIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range address {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

// consul 服务端会自己发送请求，来进行健康检查
func consulCheck(w http.ResponseWriter, r *http.Request) {
	s := time.Now().Format("2006-01-02 15:04:05") + " consul check [" + fmt.Sprint(count) + "]" + " remote:" + r.RemoteAddr + r.URL.String()
	fmt.Println(s)
	_, _ = fmt.Fprintln(w, s)
	count++
}

func runServer() {
	cfg := consulapi.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	var registration consulapi.AgentServiceRegistration
	registration.ID = "serverNode_1"       // 服务节点的名称
	registration.Name = "serverNode"       // 服务名称
	registration.Port = 9527               // 服务端口
	registration.Tags = []string{"v1"}     // tag，可以为空
	registration.Address = "127.0.0.1"     // 服务 IP
	check := &consulapi.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", // check失败后30秒删除本服务，注销时间，相当于过期时间
		// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，registration 会传到 Health.Check 函数中
	}
	registration.Check = check

	err = client.Agent().ServiceRegister(&registration)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

	http.HandleFunc("/check", consulCheck)
	err = http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)
	if err != nil {
		log.Fatal("register server error : ", err)
	}
}

func main() {
	runServer()
}
