package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	gcmd "github.com/HellseggaMosken/go-cmd"
)

func main() {
	cmd := gcmd.New("sfs", "Start a simple file server.").
		Flag(gcmd.FlagTypeValue,
			"p", "port",
			"Set the server port, default is 2579").
		Service(func(ctx *gcmd.Context) error {
			p := "2579"
			if v, ok := ctx.Long("port"); ok {
				p = v.(string)
			}
			if _, err := strconv.Atoi(p); err != nil {
				return fmt.Errorf("port should be a number, but got '%v'", p)
			}
			return startServer(ctx.Working(), p)
		})

	err := cmd.RunWithArgs()
	if err != nil {
		fmt.Println(err)
	}
}

func startServer(dir string, port string) error {
	fmt.Println("serving \"" + dir + "\" at:")
	fmt.Println("http://localhost:" + port)
	if ipv4, err := getLocalIPv4(); err == nil {
		fmt.Println("http://" + ipv4 + ":" + port)
	} else {
		fmt.Println(err)
	}

	http.Handle("/", http.FileServer(http.Dir(dir)))
	return http.ListenAndServe(":"+port, nil)
}

func getLocalIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, isIPNet := addr.(*net.IPNet)
		if !isIPNet {
			continue
		}

		if !ipNet.IP.IsLoopback() && ipNet.IP.IsGlobalUnicast() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", errors.New("can't get local ipv4 ip")
}
