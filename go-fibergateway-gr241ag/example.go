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

	// list example
	servers, err := c.VirtualServers.List()
	if err != nil {
		fmt.Println("list virtual servers failed:", err.Error())
		os.Exit(1)
	}
	natsFound := false
	for _, server := range servers {
		fmt.Println("virtual server:", server)
		if !natsFound && server.ServerName == "NATS" {
			natsFound = true
		}
	}

	// delete example
	if natsFound {
		fmt.Println("deleting virtual server")
		err = c.VirtualServers.Delete(client.VirtualServerDeleteInput{
			ExternalPortStart: "4222",
			InternalPortStart: "4222",
			Protocol:          "TCP",
			ServerIPAddress:   "192.168.1.2",
		})
		if err != nil {
			fmt.Println("virtual server creation failed:", err.Error())
			os.Exit(1)
		}
	}

	// create example
	fmt.Println("creating virtual server")
	err = c.VirtualServers.Create(client.VirtualServerCreateInput{
		ExternalPortStart: "4222",
		InternalPortStart: "4222",
		Protocol:          "TCP",
		ServerName:        "NATS",
		ServerIPAddress:   "192.168.1.2",
		WANInterface:      "veip0.1",
	})
	if err != nil {
		fmt.Println("virtual server creation failed:", err.Error())
		os.Exit(1)
	}

	s, err := c.VirtualServers.Read(client.VirtualServerReadInput{
		ExternalPortStart: "4222",
		InternalPortStart: "4222",
		Protocol:          "TCP",
		ServerIPAddress:   "192.168.1.2",
	})
	if err != nil {
		fmt.Println("virtual server read failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("the created virtual server is", s)
}
