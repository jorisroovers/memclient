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

func exec_command(server string, command string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		// TODO(jorisroovers): handle error
	}
	fmt.Fprintf(conn, command + "\r\n")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" || line == "STORED" {
			break
		}
		fmt.Println(line)
	}

	if err != nil {
		// TODO(jorisroovers): handle error
	}
	fmt.Fprintf(conn, "quit\r\n")
}

func set(server string, key string, value string) {
	flags := "0" // TODO(jorisroovers): support flags
	expiration := 0 // 0 = unlimited
	command := fmt.Sprintf("set %s %s %d %d\r\n%s", key, flags, expiration, len(value), value)
	exec_command(server, command)
}

func get(server string, key string) {
	command := fmt.Sprintf("get %s\r\n", key)
	exec_command(server, command)
}

func main() {
	cp := cli.App("memclient", "Simple command-line client for Memcached")
	host := cp.StringOpt("host h", "localhost", "Memcached host (or IP)")
	port := cp.StringOpt("port p", "11211", "Memcached port")

	cp.Command("set", "Sets a key value pair", func(cmd *cli.Cmd) {
		key := cmd.StringArg("KEY", "", "Key to set a value for")
		value := cmd.StringArg("VALUE", "", "Value to set")
		cmd.Action = func() {
			server := *host + ":" + *port
			set(server, *key, *value)
		}
	})
	cp.Command("get", "Retrieves a key", func(cmd *cli.Cmd) {
		key := cmd.StringArg("KEY", "", "Key to set a value for")
		cmd.Action = func() {
			server := *host + ":" + *port
			get(server, *key)
		}
	})
	cp.Command("version", "Prints the version", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Println("memclient, version", VERSION)
		}
	})

	cp.Run(os.Args)
}
