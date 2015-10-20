# memclient #
[![Build Status](https://travis-ci.org/jorisroovers/memclient.svg?branch=master)]
(https://travis-ci.org/jorisroovers/memclient)
[![GoDoc](https://godoc.org/github.com/jorisroovers/memclient?status.svg)]
(https://godoc.org/github.com/jorisroovers/memclient)

Memclient is a simple memcached commandline client written in Go.

**Memclient is very much under development and still missing some fundamental features**

Installation:
```
# Linux
wget https://raw.githubusercontent.com/jorisroovers/memclient/bin/linux/memclient
# Mac OS X
wget https://raw.githubusercontent.com/jorisroovers/memclient/bin/osx/memclient
```
For the latest functionality, build memclient from source:
```
go build memclient.go
```

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
  list         Lists all keys

Run 'memclient COMMAND --help' for more information on a command.
```

# Why memclient? #
1. I couldn't find a simple commandline tool for easily inspecting the data stored by a memcached server
2. I was looking for a good opportunity to get some more experience with Go :-)

# Contributing #
I'd love for you to contribute to memclient. Just open a pull request and I'll get right on it! 
You can find a wishlist below, but I'm obviously open to any suggestions you might have!

# Development #
There is a Vagrantfile in this repository that can be used for development.

```bash  
vagrant up
vagrant ssh
```

To run/build memclient:

```bash
go get github.com/jawher/mow.cli
go run memclient.go set foo bar
go run memclient.go get foo
go build memclient.go
```

To run tests:
```
go test
```

Memclient uses the [Memcached protocol](https://github.com/memcached/memcached/blob/master/doc/protocol.txt) to
talk to a memcached server.

# Wishlist #
- Pretty command output
- Better error handling
- Support for other memcached commands
- Support for additional command options
- Use GoDep for dependency management
- Have a look at [GoDownDoc](https://github.com/robertkrimen/godocdown) for README generation from ```doc.go```
- ...