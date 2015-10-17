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

type CommandExecuter interface {
	execute(command string) []string
	Close()
}

type MemcachedCommandExecuter struct {
	connection net.Conn
}

type  memClient struct {
	server   string
	executer CommandExecuter
}

func MemClient(server string) (*memClient, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	return &memClient{
		server: server,
		executer: &MemcachedCommandExecuter{
			connection: conn,
		},
	}, nil
}

func (client *memClient) Close() {
	if client.executer != nil {
		client.executer.Close()
	}
}

func (executer *MemcachedCommandExecuter) execute(command string) []string {
	fmt.Fprintf(executer.connection, command + "\r\n")
	scanner := bufio.NewScanner(executer.connection)
	var result []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" || line == "STORED" {
			break
		}
		result = append(result, line)
	}
	fmt.Fprintf(executer.connection, "quit\r\n")
	return result
}

func (executer *MemcachedCommandExecuter) Close() {
	executer.connection.Close()
}


func (client *memClient) Set(key string, value string) {
	flags := "0" // TODO(jorisroovers): support flags
	expiration := 0 // 0 = unlimited
	command := fmt.Sprintf("set %s %s %d %d\r\n%s", key, flags, expiration, len(value), value)
	client.executer.execute(command)
}

func (client *memClient) Get(key string) string {
	command := fmt.Sprintf("get %s\r\n", key)
	result := client.executer.execute(command)
	return result[1]
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
			fmt.Println(client.Get(*key))
			client.Close()
		}
	})

	cp.Run(os.Args)
}
