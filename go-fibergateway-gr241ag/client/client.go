package client

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Cristofori/kmud/telnet"
)

var ErrorNotFound = errors.New("resource not found")

type Client struct {
	host    string
	options ConnectOptions

	conn   *net.TCPConn
	telnet *telnet.Telnet

	VirtualServers VirtualServers
}

type ConnectOptions struct {
	LoginPrompt    string
	CommandsPrompt string
	NetworkType    string
	Password       string
	PasswordPrompt string
	Port           string
	Username       string
}

var DefaultOptions = ConnectOptions{
	CommandsPrompt: "/cli>",
	LoginPrompt:    "Login:",
	NetworkType:    "tcp",
	PasswordPrompt: "Password:",
	Port:           "23",
}

func Connect(host string, options ConnectOptions) (Client, error) {
	// handle default options
	if options.Port == "" {
		options.Port = DefaultOptions.Port
	}
	if options.LoginPrompt == "" {
		options.LoginPrompt = DefaultOptions.LoginPrompt
	}
	if options.PasswordPrompt == "" {
		options.PasswordPrompt = DefaultOptions.PasswordPrompt
	}
	if options.NetworkType == "" {
		options.NetworkType = DefaultOptions.NetworkType
	}
	if options.CommandsPrompt == "" {
		options.CommandsPrompt = DefaultOptions.CommandsPrompt
	}

	c := Client{
		host:    host,
		options: options,
	}

	// connect to server
	address := fmt.Sprintf("%s:%s", host, options.Port)
	tcpServer, err := net.ResolveTCPAddr(options.NetworkType, address)
	if err != nil {
		return c, fmt.Errorf("TCP addr resolution failed: %w", err)
	}

	conn, err := net.DialTCP(options.NetworkType, nil, tcpServer)
	if err != nil {
		return c, fmt.Errorf("dial failed: %w", err)
	}
	c.conn = conn
	c.telnet = telnet.NewTelnet(conn)

	// handle authentication
	if options.Username != "" && options.Password != "" {
		promptFound, data, err := c.WaitForPrompt(options.LoginPrompt)
		if err != nil {
			return c, fmt.Errorf("failed to find login prompt: %w - telnet output is %s", err, string(data))
		}
		if !promptFound {
			return c, fmt.Errorf("failed to find login prompt - telnet output is %s", string(data))
		}
		err = c.WriteTelnet(options.Username)
		if err != nil {
			return c, fmt.Errorf("login input failed: %w", err)
		}

		promptFound, data, err = c.WaitForPrompt(options.PasswordPrompt)
		if err != nil {
			return c, fmt.Errorf("failed to find password prompt: %w - telnet output is %s", err, string(data))
		}
		if !promptFound {
			return c, fmt.Errorf("failed to find password prompt - telnet output is %s", string(data))
		}
		err = c.WriteTelnet(options.Password)
		if err != nil {
			return c, fmt.Errorf("password input failed: %w", err)
		}
	}

	// wait for commands prompt
	promptFound, data, err := c.WaitForPrompt(options.CommandsPrompt)
	if err != nil {
		return c, fmt.Errorf("failed to find commands prompt: %w - telnet output is %s", err, string(data))
	}
	if !promptFound {
		return c, fmt.Errorf("failed to find commands prompt - telnet output is %s", string(data))
	}

	// init clients
	c.VirtualServers = virtualServers{client: c}

	return c, nil
}

func (c Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c Client) WriteTelnet(value string) error {
	_, err := c.telnet.Write([]byte(fmt.Sprintf("%s\n", value)))

	return err
}

func (c Client) WaitForPrompt(prompt string) (bool, []byte, error) {
	b := bytes.NewBuffer([]byte{})
	promptFound := false
	tries := 5
	for {
		data := make([]byte, 1024)
		c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := c.telnet.Read(data)
		if err != nil {
			return false, b.Bytes(), fmt.Errorf("failed waiting for prompt: %w", err)
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
