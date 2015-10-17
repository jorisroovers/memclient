package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"regexp"
	"strconv"
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
	if len(result) == 2 {
		return result[1]
	}
	return "[NOT FOUND]"
}


func (client *memClient) ListKeys() []string {
	keys := []string{}
	result := client.executer.execute("stats items\r\n")

	// identify all slabs and their number of items by parsing the 'stats items' command
	r, _ := regexp.Compile("STAT items:([0-9]*):number ([0-9]*)")
	slabCounts := map[int]int{}
	for _, stat := range result {
		matches := r.FindStringSubmatch(stat)
		if len(matches) == 3 {
			slabId, _ := strconv.Atoi(matches[1])
			slabItemCount, _ := strconv.Atoi(matches[2])
			slabCounts[slabId] = slabItemCount
		}
	}

	// For each slab, dump all items and add each key to the `keys` slice
	r, _ = regexp.Compile("ITEM (.*?) .*")
	for slabId, slabCount := range slabCounts {
		command := fmt.Sprintf("stats cachedump %v %v\n", slabId, slabCount)
		commandResult := client.executer.execute(command)
		for _, item := range commandResult {
			matches := r.FindStringSubmatch(item)
			keys = append(keys, matches[1])
		}
	}

	return keys
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
	cp.Command("list", "Lists all keys", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client := createClient(host, port)
			keys := client.ListKeys()
			for _, item := range keys {
				fmt.Println(item)
			}
			client.Close()
		}
	})

	cp.Run(os.Args)
}
