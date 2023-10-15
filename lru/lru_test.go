package lru

import (
	"fmt"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := NewLRU[string, String](180, nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		fmt.Println(v)
		t.Fatalf("cache hit key1=1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "val1", "val2", "val3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewLRU[string, String](cap, nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	fmt.Println(lru.Len())
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest key1 failed")
	}
}

func TestCache_OnEvicted(t *testing.T) {
	evictCounter := 0
	callback := func(k string, v String) {
		fmt.Println(k, v)
		evictCounter++
	}
	lru := NewLRU[string, String](10, callback)
	lru.Add("kl", String("vl"))
	lru.Add("k2", String("v2"))
	lru.Add("k3", String("v3"))
	lru.Add("k4", String("v4"))
	fmt.Println(evictCounter)
	fmt.Println(lru.Len())

}

//
//func TestCache_Add(t *testing.T) {
//	lru := New(int64(0), nil)
//	lru.Add("k1", String("v1"))
//	lru.Add("k2", String("v2"))
//
//	if lru.nbytes != int64(len("k1")+len("v1")+len("k2")+len("v2")) {
//		t.Fatal("expected 4 but got", lru.nbytes)
//	}
//}

// docs https://github.com/hashicorp/golang-lru/blob/main/simplelru/lru_test.go
