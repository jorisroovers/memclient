// +build integration

package main
import (
	"reflect"
	"sort"
	"testing"
)


const (
	server = "localhost:11211"
)

// Actual integration tests

func TestSetListGetDeleteStat(t *testing.T) {
	memClient, err := MemClient(server)
	if err != nil {
		t.Errorf("Failed to connect to the memcached server")
	}
	// assert that we start with an empty server
	keys := memClient.ListKeys()
	if !reflect.DeepEqual([]string{}, keys) {
		t.Errorf("Expected %v, got %v. Make sure the memcached is empty", []string{}, keys)
	}

	// Set some values and assert that they've been set correctly
	memClient.Set("foo", "bar", 0)
	memClient.Set("test1", "value1", 0)
	memClient.Set("test2", "value2", 0)

	keys = memClient.ListKeys()
	sort.Strings(keys)
	expected := []string{"foo", "test1", "test2"}
	if !reflect.DeepEqual(expected, keys) {
		t.Errorf("Expected %v, got %v.", expected, keys)
	}

	val, ok := memClient.Get("foo")
	if !ok || !reflect.DeepEqual([]string{"bar"}, val) {
		t.Errorf("Expected %v, got %v.", []string{"bar"}, val)
	}
	val, ok = memClient.Get("test1")
	if !ok || !reflect.DeepEqual([]string{"value1"}, val) {
		t.Errorf("Expected %v, got %v.", []string{"value1"}, val)
	}
	val, ok = memClient.Get("test2")
	if !ok || !reflect.DeepEqual([]string{"value2"}, val) {
		t.Errorf("Expected %v, got %v.", []string{"value2"}, val)
	}

	// Get a value that doesn't exist
	_, ok = memClient.Get("foobar")
	if ok {
		t.Errorf("memClient.Get('foobar') is returning a value, but it shouldn't")
	}

	// Delete a value and assert it can no longer be found
	memClient.Delete("foo")
	keys = memClient.ListKeys()
	sort.Strings(keys)
	expected = []string{"test1", "test2"}
	if !reflect.DeepEqual(expected, keys) {
		t.Errorf("Expected %v, got %v.", expected, keys)
	}

	// Try to get a value that we didn't set and assert it cannot be found
	_, ok = memClient.Get("foo")
	if ok {
		t.Errorf("memClient.Get('foobar') is returning a value, but it shouldn't")
	}

	// Get the cmd_get statistic and verify it is set to 5 (since we did 5 get operations)
	cmdGet, ok := memClient.Stat("cmd_get")
	expectedCmdGet := Stat{"cmd_get", "5"}
	if !ok || cmdGet != expectedCmdGet {
		t.Errorf("cmd_get statistic is not 5")
	}

}

func TestVersion(t *testing.T) {
	memClient, err := MemClient(server)
	if err != nil {
		t.Errorf("Failed to connect to the memcached server")
	}
	version := memClient.Version()
	expected := "VERSION 1.4.14 (Ubuntu)"
	if version != expected {
		t.Errorf("Incorrect version string, expected '%v', got '%v'", expected, version)
	}

}