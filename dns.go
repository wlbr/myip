package main

import (
	"fmt"
	"net"
)

func main() {

	// host, _ := net.LookupAddr("193.99.144.85")
	host, _ := net.LookupAddr("5.158.147.174")
	fmt.Println(host)

	addresses, _ := net.LookupHost(host[0])
	fmt.Println(addresses)
}
