package storage

import (
		"testing"
)

func TestSet(t *testing.T) {
	s := NewStorage()
	s.Set("key1", "value1")
}

func TestGet(t *testing.T) {
	s := NewStorage()
	s.Set("key1", "value1")

	value, ok := s.Get("key1")
	if ok {
		t.Errorf("Get failed")
	}
	if value != "value1" {
		t.Errorf("Get returned wrong value: got %v want value1", value)
	}
}