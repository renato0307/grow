package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Cristofori/kmud/telnet"
)

const (
	HOST = "192.168.1.254"
	PORT = "23"
	TYPE = "tcp"
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

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	telnet := telnet.NewTelnet(conn)

	loginFound, _, err := waitForPrompt(telnet, "Login:")
	if err != nil {
		println("Wait for login failed:", err.Error())
		os.Exit(1)
	}
	if !loginFound {
		println("No login")
		os.Exit(1)
	}

	_, err = telnet.Write([]byte("meo\n"))
	if err != nil {
		println("Login input failed:", err.Error())
		os.Exit(1)
	}

	passwordFound, _, err := waitForPrompt(telnet, "Password:")
	if err != nil {
		println("Wait for password failed:", err.Error())
		os.Exit(1)
	}
	if !passwordFound {
		println("No password")
		os.Exit(1)
	}

	_, err = telnet.Write([]byte(fmt.Sprintf("%s\n", os.Getenv("ROUTER_PASSWORD"))))
	if err != nil {
		println("Password input failed:", err.Error())
		os.Exit(1)
	}

	cliPromptFound, _, err := waitForPrompt(telnet, "/cli>")
	if err != nil {
		println("Wait for cli prompt failed:", err.Error())
		os.Exit(1)
	}
	if !cliPromptFound {
		println("No cli prompt")
		os.Exit(1)
	}

	_, err = telnet.Write([]byte("nat/virtual-servers/show\n"))
	if err != nil {
		println("Command input failed:", err.Error())
		os.Exit(1)
	}

	cliPromptFound, data, err := waitForPrompt(telnet, "/cli>")
	if err != nil {
		println("Wait for virtual servers show failed:", err.Error())
		os.Exit(1)
	}
	if !cliPromptFound {
		println("No cli prompt")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	servers := []VirtualServer{}
	var server VirtualServer
	newServer := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "/cli>") {
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

	for _, server := range servers {
		fmt.Println("virtual server:", server)
	}

}

func waitForPrompt(telnet *telnet.Telnet, prompt string) (bool, []byte, error) {
	b := bytes.NewBuffer([]byte{})
	promptFound := false
	tries := 5
	for {
		data := make([]byte, 1024)
		n, err := telnet.Read(data)
		if err != nil {
			return false, nil, fmt.Errorf("read failed: %w", err)
		}
		b.Write(data[:n])

		promptFound = strings.HasSuffix(strings.TrimSpace(b.String()), prompt)
		if promptFound {
			break
		}
		tries--
		if tries == 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	return promptFound, b.Bytes(), nil
}
