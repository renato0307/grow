package client

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type VirtualServer struct {
	ServerName        string
	ExternalPortStart string
	ExternalPortEnd   string
	Protocol          string
	InternalPortStart string
	InternalPortEnd   string
	ServerIPAddress   string
	WANInterface      string
	Origin            string
}

type VirtualServers interface {
	Create(VirtualServer) error
	Delete(name string) error
	List() ([]VirtualServer, error)
}

type virtualServers struct {
	client Client
}

func (vs virtualServers) Create(VirtualServer) error {
	return nil
}

func (vs virtualServers) Delete(name string) error {
	return errors.New("not implemented")
}

func (vs virtualServers) List() ([]VirtualServer, error) {
	err := vs.client.WriteTelnet("nat/virtual-servers/show")
	if err != nil {
		return nil, fmt.Errorf("failed to list virtual servers: %w", err)
	}

	promptFound, data, err := vs.client.WaitForPrompt(vs.client.options.CommandsPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to find commands prompt: %w", err)
	}
	if !promptFound {
		return nil, fmt.Errorf("failed to find commands prompt")
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	servers := []VirtualServer{}
	var server VirtualServer
	newServer := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, vs.client.options.CommandsPrompt) {
			break
		}
		if strings.Contains(line, "---------------------") {
			newServer = !newServer
			if !newServer {
				servers = append(servers, server)
				server = VirtualServer{}
				newServer = true
			}
			continue
		}

		split := strings.Split(line, ":")
		if len(split) < 2 {
			continue // invalid field
		}

		field := strings.ReplaceAll(strings.Trim(split[0], ": "), " ", "")
		value := strings.Trim(split[1], ": ")

		switch field {
		case "ExternalPortEnd":
			server.ExternalPortEnd = value
		case "ExternalPortStart":
			server.ExternalPortStart = value
		case "InternalPortEnd":
			server.InternalPortEnd = value
		case "InternalPortStart":
			server.InternalPortStart = value
		case "Origin":
			server.Origin = value
		case "Protocol":
			server.Protocol = value
		case "ServerIPAddress":
			server.ServerIPAddress = value
		case "ServerName":
			server.ServerName = value
		case "WANInterface":
			server.WANInterface = value
		}
	}

	return servers, nil
}
