# memclient with Apple Silicon Support

Memclient is a simple memcached commandline client written in Go. Originally from: https://github.com/jorisroovers/memclient

This version differs from the original in that it
1. Produces a release which works with Apple Silicon
2. Removes the Vagrant/Travis stuff from code since I don't use that
3. Now uses go modules
4. Specify `cli` in the `import` so golang doesn't remove it automatically
5. Updated gitignore

**Memclient is very much under development and still missing some fundamental features**

Example usage:

```bash
# To set a key:
memclient set mykey myvalue
# To retrieve a key:
memclient get mykey
# To delete a key
memclient delete mykey
# To list all available keys
memclient list
# Flush all keys (they will still show in 'list', but will return 'NOT FOUND' when fetched using 'memclient get')
memclient flush
# Print the server version
memclient version

```

Other commands:

```bash
Usage: memclient [OPTIONS] COMMAND [arg...]

Simple command-line client for Memcached

Options:
  -v, --version=false      Show the version and exit
  --host, -h="localhost"   Memcached host (or IP)
  --port, -p="11211"       Memcached port

Commands:
  set          Sets a key value pair
  get          Retrieves a key
  delete       Deletes a key
  flush        Flush all cache keys (they will still show in 'list', but will return 'NOT FOUND')
  version      Show server version
  list         Lists all keys
  stats        Print server statistics
  stat         Print a specific server statistic

Run 'memclient COMMAND --help' for more information on a command.
```

# Contributing

Memclient uses the [Memcached protocol](https://github.com/memcached/memcached/blob/master/doc/protocol.txt) to
talk to a memcached server.

To build memclient from source:

```
go build memclient.go
```