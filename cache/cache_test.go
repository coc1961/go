package cache

import (
	"fmt"
	"testing"
)

const (
	KEY1 = "key1"
	KEY2 = "key2"
	KEY3 = "key3"
	KEY4 = "key4"
	KEY5 = "key5"
)

func TestSetGet(t *testing.T) {
	cache := NewCache("127.0.0.1:11211")

	cache.DeleteAll()

	err := cache.Set(KEY1, []byte("Esta es una Prueba"))
	printError(t, err)

	it, err := cache.Get(KEY1)
	printError(t, err)

	fmt.Println("TestSetGet =", string(it.Value), "(Ok)")
}

func printError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
