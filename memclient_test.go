// +build !integration

package main

import (
	"reflect"
	"testing"
)

// Test Helpers
type MockedCommandExecuter struct {
	t                *testing.T
	executedCommands []string
	returnValues     map[string][]string
	closed           bool
}

func (executer *MockedCommandExecuter) execute(command string, responseDelimiters []string) []string {
	executer.executedCommands = append(executer.executedCommands, command)
	returnVal, ok := executer.returnValues[command]
	if ok {
		return returnVal
	}
	return []string{}
}

func (executer *MockedCommandExecuter) Close() {
}

func (executer *MockedCommandExecuter) addReturnValue(command string, returnValue []string) {
	executer.returnValues[command] = returnValue
}

/*
	Asserts that a given slice of commands have been called executed against the command executer.
 */
func (executer *MockedCommandExecuter) assertCommands(expectedCommands []string) {
	if !reflect.DeepEqual(executer.executedCommands, expectedCommands) {
		executer.t.Errorf("Executed command were '%v', expected '%v'", executer.executedCommands, expectedCommands)
	}
}

func createTestClient(t *testing.T) (*memClient, *MockedCommandExecuter) {
	executer := &MockedCommandExecuter{t, []string{}, map[string][]string{}, false}
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
	executer.assertCommands([]string{"get testkey\r\n"})
}

func TestSet(t *testing.T) {
	client, executer := createTestClient(t)
	client.Set("testkey", "testval", 123)
	executer.assertCommands([]string{"set testkey 0 123 7\r\ntestval\r\n"})
}

func TestDelete(t *testing.T) {
	client, executer := createTestClient(t)
	client.Delete("testkey")
	executer.assertCommands([]string{"delete testkey\r\n"})
}

func TestVersion(t *testing.T) {
	client, executer := createTestClient(t)
	version := client.Version()
	executer.assertCommands([]string{"version \r\n"})
	if (version != "UNKNOWN") {
		t.Errorf("Received version does not match expected version (%v!=%v)", version, "VERSION myversion.1234")
	}

	executer.addReturnValue("version \r\n", []string{"VERSION myversion.1234"})
	version = client.Version()
	if (version != "VERSION myversion.1234") {
		t.Errorf("Received version does not match expected version (%v!=%v)", version, "VERSION myversion.1234")
	}
}

func TestFlush(t *testing.T) {
	client, executer := createTestClient(t)
	client.Flush()
	executer.assertCommands([]string{"flush_all \r\n"})
}

func TestListKeys(t *testing.T) {
	// setup testcase
	client, executer := createTestClient(t)
	executer.addReturnValue("stats items\r\n", []string{"STAT items:123:number 456"})
	executer.addReturnValue("stats cachedump 123 456\n", []string{"ITEM foobar ignored", "ITEM testkey ignored"})

	keys := client.ListKeys()

	// validate that the result is correct and that the expected commands where executed
	expectedKeys := []string{"foobar", "testkey"}
	if (!reflect.DeepEqual(keys, expectedKeys)) {
		t.Errorf("Returned cache keys incorrect (%v!=%v)", keys, expectedKeys)
	}
	executer.assertCommands([]string{"stats items\r\n", "stats cachedump 123 456\n"})

}
