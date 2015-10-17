package main

import (
	"testing"
)

// Test Helpers
type MockedCommandExecuter struct {
	t               *testing.T
	executedCommand string
}

func (executer *MockedCommandExecuter) execute(command string) {
	executer.executedCommand = command
}

func (executer *MockedCommandExecuter) assertCommand(expectedCommand string) {
	if executer.executedCommand != expectedCommand {
		executer.t.Errorf("Executed command was '%v', expected '%v'", executer.executedCommand, expectedCommand)
	}
}

func createTestClient(t *testing.T) (*memClient, *MockedCommandExecuter) {
	executer := &MockedCommandExecuter{t, ""}
	client := &memClient{
		server: "foo",
		executer: executer,
	}
	return client, executer
}

// Actual tests

func TestMemClient(t *testing.T) {
	_, err := MemClient("foo:1234")
	if err == nil {
		t.Errorf("Memclient should return an error for foo:1234")
	}
}

func TestGet(t *testing.T) {
	client, executer := createTestClient(t)
	client.Get("testkey")
	executer.assertCommand("get testkey\r\n")
}

func TestSet(t *testing.T) {
	client, executer := createTestClient(t)
	client.Set("testkey", "testval")
	executer.assertCommand("set testkey 0 0 7\r\ntestval")
}
