package main

import (
	"testing"
	"net"
)

type DummyConnection struct {
	net.Conn
}


func TestMemClient(t *testing.T) {
	_, err := MemClient("foo:1234")
	if err == nil {
		t.Errorf("Memclient should return an error for foo:1234")
	}
	//	client.connection = &DummyConnection{}
}