package main

import (
	"fmt"
	"os"

	"github.com/renato0307/grow/go-fibergateway-gr241ag/client"
)

func main() {
	c, err := client.Connect("192.168.1.254", client.ConnectOptions{
		Username: "meo",
		Password: os.Getenv("ROUTER_PASSWORD"),
	})
	if err != nil {
		fmt.Println("connection to router failed:", err.Error())
		os.Exit(1)
	}
	defer c.Close()
	fmt.Println("connected to router")

	servers, err := c.VirtualServers.List()
	if err != nil {
		fmt.Println("list virtual servers failed:", err.Error())
		os.Exit(1)
	}
	for _, server := range servers {
		fmt.Println("virtual server:", server)
	}
}
