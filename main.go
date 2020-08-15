package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
)

// 获取本机网卡IPv4地址
func getLocalIPv4() (string, error) {
	// 获取所有网卡
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	// 取第一个非lo的网卡IP
	for _, addr := range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIPNet := addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("can't get local ipv4 ip")
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic("can't get the current working dir's path")
	}

	port := "8080"

	fmt.Println("server at:")
	fmt.Println("http://localhost:" + port)
	if ipv4, err := getLocalIPv4(); err == nil {
		fmt.Println("http://" + ipv4 + ":" + port)
	} else {
		fmt.Println(err)
	}

	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(":"+port, nil)
}
