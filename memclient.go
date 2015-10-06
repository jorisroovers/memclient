package main

import (
	"fmt"
	"flag"
	"net"
	"bufio"
	"os"
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

func assert_arg_length(length int) {
	args := flag.Args()
	if len(args) < length {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: memclient [options] command\n")
		fmt.Println("commands:")
		fmt.Println("  set [key] [value]")
		fmt.Println("    Sets a key value pair")
		fmt.Println("  get [key]")
		fmt.Println("    Retrieves a key")
		fmt.Println("options:")
		flag.PrintDefaults()
	}

	server := flag.String("server", "localhost:11211", "Memcached server:port")
	//	exec_command(*server, "stats")
	flag.Parse()
	args := flag.Args()
	assert_arg_length(1)

	command := args[0]
	switch command {
	case "set":
		assert_arg_length(3)
		set(*server, args[1], args[2])
	case "get":
		assert_arg_length(2)
		get(*server, args[1])
	}
}
