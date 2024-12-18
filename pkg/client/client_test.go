package client

import (
	"fmt"
	"log"
	"testing"
	"time"
)


func TestClient(t *testing.T) {

	// Create a client
	clients := NewClient([]string{
		"localhost:8001",
	})
	

	// Test basic operations
	t.Log("Testing basic operations...")

	// Test SET
	err := clients.Set("name", "John Doe")
	if err != nil {
		log.Fatal("SET failed:", err)
	}
	fmt.Println("Successfully set key 'name'")

	time.Sleep(2 * time.Second)

	// Test GET
	value, exists, err := clients.Get("name")
	if err != nil {
		log.Fatal("GET failed:", err)
	}
	if exists {
		fmt.Printf("Successfully retrieved value: %s\n", value)
	} else {
		fmt.Println("Key not found")
	}

	// Test non-existent key
	// value, exists, err = clients.Get("nonexistent")
	// if err != nil {
	// 	log.Fatal("GET failed:", err)
	// }
	// if !exists {
	// 	fmt.Println("Correctly reported non-existent key")
	// }
	time.Sleep(5 * time.Second)
	// Test DELETE
	err = clients.Delete("name")
	if err != nil {
		log.Fatal("DELETE failed:", err)
	}
	fmt.Println("Successfully deleted key 'name'")

	time.Sleep(5 * time.Second)
	// Verify deletion
	_, exists, err = clients.Get("name")
	if err != nil {
		log.Fatal("GET failed:", err)
	}

	if !exists {
		fmt.Println("Successfully verified deletion")
	}

	time.Sleep(5 * time.Second)
	// Test fault tolerance
	fmt.Println("\nTesting fault tolerance...")

	// Set a value
	err = clients.Set("test-key", "test-value")
	if err != nil {
		log.Fatal("SET failed:", err)
	}

	// Try to get the value from different nodes by using different client instances
	client_kvdb := []*Client{
		NewClient([]string{"localhost:8001"}),
		NewClient([]string{"localhost:8002"}),
		NewClient([]string{"localhost:8003"}),
	}

	for i, cli := range client_kvdb {
		value, exists, err := cli.Get("test-key")
		time.Sleep(2 * time.Second)
		if err != nil {
			fmt.Printf("Node %d: Error getting value: %v\n", i+1, err)
		} else if exists {
			fmt.Printf("Node %d: Successfully retrieved value: %s\n", i+1, value)
		} else {
			fmt.Printf("Node %d: Key not found\n", i+1)
		}
	}
}
