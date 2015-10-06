# memclient #

Memclient is a simple memcached commandline client written in Go.

** Memclient is very much under development and still missing some fundamental features **

Installation:
```
# Linux
wget https://raw.githubusercontent.com/jorisroovers/memclient/bin/linux/memclient
# Mac OS X
wget https://raw.githubusercontent.com/jorisroovers/memclient/bin/osx/memclient
```

Example usage:
```bash
# To set a key:
memclient set mykey myvalye
# To retrieve a key:
memclient get mykey
```

Other commands:
```bash
Usage: memclient [options] command
commands:
  set [key] [value]
    Sets a key value pair
  get [key]
    Retrieves a key
options:
  -server string
    	Memcached server:port (default "localhost:11211")
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
go run memclient.go set foo bar
go run memclient.go get foo
go build memclient.go
```

Memclient uses the [Memcached protocol](https://github.com/memcached/memcached/blob/master/doc/protocol.txt) to
talk to a memcached server.

# Wishlist #
- Listing all keys
- Pretty command output
- Better error handling
- Unit tests!
- Support for other memcached commands
- ...