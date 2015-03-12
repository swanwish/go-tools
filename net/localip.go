package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ips, err := GetLocalIPAddrs()
	if err != nil {
		log.Println("Failed to get ip address.", err)
		return
	}
	for _, ip := range ips {
		fmt.Println(ip)
	}
}

func GetLocalIPAddrs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	ips := make([]string, 0)
	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}

		}
	}
	return ips, nil
}
