package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"github.com/jawher/mow.cli"
)

const (
	VERSION = "0.1.0dev"
)

type  memClient struct {
	server     string
	connection net.Conn
}

func MemClient(server string) (*memClient, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	return &memClient{
		server: server,
		connection: conn,
	}, nil
}

func (client *memClient) Close() {
	if client != nil {
		client.Close()
	}
}

/*
 * Executes a memcached command and returns the server output.
 */
func (client *memClient) exec_command(command string) {
	fmt.Fprintf(client.connection, command + "\r\n")
	scanner := bufio.NewScanner(client.connection)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" || line == "STORED" {
			break
		}
		fmt.Println(line)
	}

	fmt.Fprintf(client.connection, "quit\r\n")
}

func (client *memClient) Set(key string, value string) {
	flags := "0" // TODO(jorisroovers): support flags
	expiration := 0 // 0 = unlimited
	command := fmt.Sprintf("set %s %s %d %d\r\n%s", key, flags, expiration, len(value), value)
	client.exec_command(command)
}

func (client *memClient) Get(key string) {
	command := fmt.Sprintf("get %s\r\n", key)
	client.exec_command(command)
}

func createClient(host, port *string) (*memClient) {
	server := *host + ":" + *port
	client, err := MemClient(server)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to", server)
		os.Exit(1)
	}
	return client
}

func main() {
	cp := cli.App("memclient", "Simple command-line client for Memcached")
	cp.Version("v version", "memclient " + VERSION)
	host := cp.StringOpt("host h", "localhost", "Memcached host (or IP)")
	port := cp.StringOpt("port p", "11211", "Memcached port")

	cp.Command("set", "Sets a key value pair", func(cmd *cli.Cmd) {
		key := cmd.StringArg("KEY", "", "Key to set a value for")
		value := cmd.StringArg("VALUE", "", "Value to set")
		cmd.Action = func() {
			client := createClient(host, port)
			client.Set(*key, *value)
			client.Close()
		}
	})
	cp.Command("get", "Retrieves a key", func(cmd *cli.Cmd) {
		key := cmd.StringArg("KEY", "", "Key to set a value for")
		cmd.Action = func() {
			client := createClient(host, port)
			client.Get(*key)
			client.Close()
		}
	})

	cp.Run(os.Args)
}
