package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/jawher/mow.cli"
)

const (
	VERSION = "0.1.0dev"
)

/*
	The CommandExecuter interface defines an entity that is able to execute memcached
	commands against a memcached server.
 */
type CommandExecuter interface {
	execute(command string, delimiters []string) []string
	Close()
}

type MemcachedCommandExecuter struct {
	connection net.Conn
}

type  memClient struct {
	server   string
	executer CommandExecuter
}

type Stat struct {
	name  string
	value string
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

func (executer *MemcachedCommandExecuter) execute(command string, responseDelimiters []string) []string {
	fmt.Fprintf(executer.connection, command)
	scanner := bufio.NewScanner(executer.connection)
	var result []string

	OUTER:
	for scanner.Scan() {
		line := scanner.Text()
		for _, delimeter := range responseDelimiters {
			if line == delimeter {
				break OUTER
			}
		}
		result = append(result, line)
		// if there is no delimiter specified, then the response is just a single line and we should return after
		// reading that first line (e.g. version command)
		if len(responseDelimiters) == 0 {
			break OUTER
		}
	}
	return result
}

func (executer *MemcachedCommandExecuter) Close() {
	executer.connection.Close()
}

/*
	Retrieves a given cache key from the memcached server.
	Returns a string array with the value and a boolean indicating whether a value was found or not.
 */
func (client *memClient) Get(key string) ([]string, bool) {
	command := fmt.Sprintf("get %s\r\n", key)
	result := client.executer.execute(command, []string{"END"})
	if len(result) >= 2 {
		// ditch the first "VALUE <key> <expiration> <length>" line
		return result[1:], true
	}
	return []string{}, false
}

/*
	Sets a given cache key on the memcached server to a given value.
 */
func (client *memClient) Set(key string, value string, expiration int) {
	flags := "0" // TODO(jorisroovers): support flags
	command := fmt.Sprintf("set %s %s %d %d\r\n%s\r\n", key, flags, expiration, len(value), value)
	client.executer.execute(command, []string{"STORED"})
}

/*
	Deletes a given cache key on the memcached server.
 */
func (client *memClient) Delete(key string) {
	command := fmt.Sprintf("delete %s\r\n", key)
	client.executer.execute(command, []string{"DELETED", "NOT_FOUND"})
}

/*
	List all cache keys on the memcached server.
 */
func (client *memClient) ListKeys() []string {
	keys := []string{}
	result := client.executer.execute("stats items\r\n", []string{"END"})

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
		commandResult := client.executer.execute(command, []string{"END"})
		for _, item := range commandResult {
			matches := r.FindStringSubmatch(item)
			keys = append(keys, matches[1])
		}
	}

	return keys
}

/*
   Get the server version.
 */
func (client *memClient) Version() string {
	result := client.executer.execute("version \r\n", []string{})
	if len(result) == 1 {
		return result[0]
	}
	return "UNKNOWN"
}

/*
	Retrieve all server statistics.
 */
func (client *memClient) Stats() []Stat {
	result := client.executer.execute("stats\r\n", []string{"END"})
	stats := []Stat{}
	for _, stat := range result {
		parts := strings.SplitN(stat, " ", 3)
		stats = append(stats, Stat{parts[1], parts[2]})
	}

	return stats
}

/*
	Retrieve a specific server statistic.
 */
func (client *memClient) Stat(statName string) (Stat, bool) {
	stats := client.Stats()
	for _, stat := range stats {
		if stat.name == statName {
			return stat, true
		}
	}
	return Stat{}, false
}

func (client *memClient) Flush() {
	client.executer.execute("flush_all \r\n", []string{"OK"})
}

/*
	Creates a memClient and deals with any errors that might occur (e.g. unable to connect to server).
 */
func createClient(host, port *string) (*memClient) {
	server := *host + ":" + *port
	client, err := MemClient(server)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to", server)
		os.Exit(1)
	}
	return client
}

/*
	This is where the magic happens.
	Creates a CLI app using mow.cli and to calls the appropriate functions on a Memclient instance according to the
	passed parameters.
 */
func main() {
	cp := cli.App("memclient", "Simple command-line client for Memcached")
	cp.Version("v version", "memclient " + VERSION)
	host := cp.StringOpt("host h", "localhost", "Memcached host (or IP)")
	port := cp.StringOpt("port p", "11211", "Memcached port")

	cp.Command("set", "Sets a key value pair", func(cmd *cli.Cmd) {
		key := cmd.StringArg("KEY", "", "Key to set a value for")
		value := cmd.StringArg("VALUE", "", "Value to set")
		expiration := cmd.IntOpt("expire", 0, "expiration time (in seconds)")
		cmd.Action = func() {
			client := createClient(host, port)
			client.Set(*key, *value, *expiration)
			client.Close()
		}
	})
	cp.Command("get", "Retrieves a key", func(cmd *cli.Cmd) {
		key := cmd.StringArg("KEY", "", "Key to set a value for")
		cmd.Action = func() {
			client := createClient(host, port)
			result, ok := client.Get(*key)
			if ok {
				for _, line := range result {
					fmt.Println(line)
				}
			} else {
				fmt.Println("[NOT FOUND]")
			}

			client.Close()
		}
	})
	cp.Command("delete", "Deletes a key", func(cmd *cli.Cmd) {
		key := cmd.StringArg("KEY", "", "Key to delete")
		cmd.Action = func() {
			client := createClient(host, port)
			client.Delete(*key)
			client.Close()
		}
	})
	description := "Flush all cache keys (they will still show in 'list', but will return 'NOT FOUND')"
	cp.Command("flush", description, func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client := createClient(host, port)
			client.Flush()
			client.Close()
		}
	})
	cp.Command("version", "Show server version", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client := createClient(host, port)
			fmt.Println(client.Version())
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
	cp.Command("stats", "Print server statistics", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client := createClient(host, port)
			stats := client.Stats()
			for _, stat := range stats {
				fmt.Printf("%-25v %v\n", stat.name, stat.value)
			}
			client.Close()
		}
	})
	cp.Command("stat", "Print a specific server statistic", func(cmd *cli.Cmd) {
		statName := cmd.StringArg("STAT", "", "Name of the statistic to get")
		cmd.Action = func() {
			client := createClient(host, port)
			stat, ok := client.Stat(*statName)
			if ok {
				fmt.Printf("%-25v %v\n", stat.name, stat.value)
			} else {
				fmt.Println("[NOT FOUND]")
			}
			client.Close()
		}
	})

	cp.Run(os.Args)
}
