/*

Memclient is a simple memcached commandline client written in Go.

**Memclient is very much under development and still missing some fundamental features**

Installation

	# Linux
	wget https://raw.githubusercontent.com/jorisroovers/memclient/bin/linux/memclient
	# Mac OS X
	wget https://raw.githubusercontent.com/jorisroovers/memclient/bin/osx/memclient


Example usage

	# To set a key:
	memclient set mykey myvalue
	# To retrieve a key:
	memclient get mykey
	# To list all available keys
	memclient list


*/
package main