package cmap

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	// Create a new map.
	m := New()

	// Sets item within map, sets "bar" under key "foo"
	m.Set("foo", "bar")

	// Retrieve item from map.
	if tmp, ok := m.Get("foo"); ok {
		bar := tmp.(string)
		fmt.Println(bar)
	}

	// Removes item under key "foo"
	m.Remove("foo")
	if tmp, ok := m.Get("foo"); ok {
		bar := tmp.(string)
		fmt.Println(bar)
	}
}
