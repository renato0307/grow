package client

import (
	"bufio"
	"fmt"
	"strings"
)

type VirtualServer struct {
	ExternalPortEnd   string
	ExternalPortStart string
	InternalPortEnd   string
	InternalPortStart string
	Origin            string
	Protocol          string
	ServerIPAddress   string
	ServerName        string
	WANInterface      string
}

type VirtualServerCreateInput VirtualServer

type VirtualServerReadInput struct {
	Name string
}

type VirtualServerDeleteInput VirtualServerReadInput

type VirtualServers interface {
	Create(VirtualServerCreateInput) error
	Delete(VirtualServerDeleteInput) error
	Read(VirtualServerReadInput) (VirtualServer, error)
	List() ([]VirtualServer, error)
}

type virtualServers struct {
	client Client
}

func (vs virtualServers) Create(server VirtualServerCreateInput) error {
	sb := strings.Builder{}
	sb.WriteString("nat/virtual-servers/create")
	sb.WriteString(fmt.Sprintf(" --ext-port-start=%s", server.ExternalPortStart))
	sb.WriteString(fmt.Sprintf(" --int-port-start=%s", server.InternalPortStart))
	sb.WriteString(fmt.Sprintf(" --protocol=%s", server.Protocol))
	sb.WriteString(fmt.Sprintf(" --server-ip=%s", server.ServerIPAddress))
	sb.WriteString(fmt.Sprintf(" --server-name=%s", server.ServerName))
	sb.WriteString(fmt.Sprintf(" --wan-intf=%s", server.WANInterface))
	if server.ExternalPortEnd != "" {
		sb.WriteString(fmt.Sprintf(" --ext-port-end=%s", server.ExternalPortEnd))
	}
	if server.InternalPortEnd != "" {
		sb.WriteString(fmt.Sprintf(" --int-port-end=%s", server.InternalPortEnd))
	}

	err := vs.client.WriteTelnet(sb.String())
	if err != nil {
		return fmt.Errorf("failed to create virtual server: %w", err)
	}

	data, err := vs.waitForCommandsPrompt()
	if err != nil {
		return err
	}

	if errorFoundOnCommandResponse(data) {
		return fmt.Errorf("failed to create virtual server: %s", data)
	}

	return nil
}

func (vs virtualServers) Delete(delServer VirtualServerDeleteInput) error {
	sb := strings.Builder{}

	server, err := vs.Read(VirtualServerReadInput(delServer))
	if err != nil {
		if err == ErrorNotFound {
			return nil
		}
		return err
	}
	sb.WriteString("nat/virtual-servers/remove")
	sb.WriteString(fmt.Sprintf(" --ext-port-start=%s", server.ExternalPortStart))
	sb.WriteString(fmt.Sprintf(" --int-port-start=%s", server.InternalPortStart))
	sb.WriteString(fmt.Sprintf(" --protocol=%s", server.Protocol))
	sb.WriteString(fmt.Sprintf(" --server-ip=%s", server.ServerIPAddress))
	if server.ExternalPortEnd != "" {
		sb.WriteString(fmt.Sprintf(" --ext-port-end=%s", server.ExternalPortEnd))
	}
	if server.InternalPortEnd != "" {
		sb.WriteString(fmt.Sprintf(" --int-port-end=%s", server.InternalPortEnd))
	}

	err = vs.client.WriteTelnet(sb.String())
	if err != nil {
		return fmt.Errorf("failed to delete virtual server: %w", err)
	}

	data, err := vs.waitForCommandsPrompt()
	if err != nil {
		return err
	}
	if errorFoundOnCommandResponse(data) {
		return fmt.Errorf("failed to delete virtual server: %s", data)
	}

	return nil
}

func (vs virtualServers) List() ([]VirtualServer, error) {
	err := vs.client.WriteTelnet("nat/virtual-servers/show")
	if err != nil {
		return nil, fmt.Errorf("failed to list virtual servers: %w", err)
	}

	data, err := vs.waitForCommandsPrompt()
	if err != nil {
		return nil, err
	}
	if errorFoundOnCommandResponse(data) {
		return nil, fmt.Errorf("failed to list virtual servers: %s", data)
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

func (vs virtualServers) Read(server VirtualServerReadInput) (VirtualServer, error) {
	servers, err := vs.List()
	if err != nil {
		return VirtualServer{}, err
	}

	for _, s := range servers {
		if server.Name == s.ServerName {
			return s, nil
		}
	}

	return VirtualServer{}, ErrorNotFound
}

func (vs virtualServers) waitForCommandsPrompt() ([]byte, error) {
	promptFound, data, err := vs.client.WaitForPrompt(vs.client.options.CommandsPrompt)
	if err != nil {
		return data, fmt.Errorf("failed to find commands prompt: %w - telnet output is %s", err, string(data))
	}
	if !promptFound {
		return data, fmt.Errorf("failed to find commands prompt - telnet output is %s", string(data))
	}
	return data, nil
}

func errorFoundOnCommandResponse(data []byte) bool {
	return strings.Contains(string(data), "ERROR") || strings.Contains(string(data), "Failed")
}
