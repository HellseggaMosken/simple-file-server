package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
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

const usage = `ACCEPTED OPTIONS:
h/help : show this usage
p/port [port number] : set server port, default 2579`

func main() {
	port := "2579"

	if len(os.Args) > 1 {
		op := os.Args[1]
		if op == "help" || op == "h" {
			fmt.Println(usage)
			return
		} else if op == "port" || op == "p" {
			if len(os.Args) > 2 {
				if _, err := strconv.Atoi(os.Args[2]); err == nil {
					port = os.Args[2]
				} else {
					fmt.Println("port should be a number, but got \"" + os.Args[2] + "\"")
					return
				}
			}
		} else {
			fmt.Println("unknown option: " + op + ", use option \"help\" for help")
			return
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		panic("can't get the current working dir's path")
	}

	fmt.Println("serving \"" + dir + "\" at:")
	fmt.Println("http://localhost:" + port)
	if ipv4, err := getLocalIPv4(); err == nil {
		fmt.Println("http://" + ipv4 + ":" + port)
	} else {
		fmt.Println(err)
	}

	http.Handle("/", http.FileServer(http.Dir(dir)))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
