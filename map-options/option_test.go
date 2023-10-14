package map_options

import (
	"fmt"
	"testing"
)

func Test_IntNewDataMapOption(t *testing.T) {
	opt := NewDataMapOption(map[string]string{
		"q":  "12312",
		"e":  "123",
		"1e": "123",
	})

	data := opt.Keys()
	fmt.Println(data)
	//
	//fmt.Println(opt.Option(1))
	//fmt.Println(opt.Options())
}

func TestName(t *testing.T) {
	data := MapKeys(map[int]int{
		1: 1,
	})
	fmt.Println(data)
}

func MapKeys[Key comparable, Val any](m map[Key]Val) []Key {
	s := make([]Key, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func BenchmarkNewDataMapOption(t *testing.B) {
	opt := NewDataMapOption(map[string]string{
		"1":  "12312",
		"23": "12312",
	})
	data := opt.Keys()
	fmt.Println(data)
}
