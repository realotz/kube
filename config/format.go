package config

import (
	"strings"
	"unsafe"
)

func format(name string) string {
	if ext := strings.Split(name, "."); len(ext) > 1 {
		return ext[len(ext)-1]
	}
	return "text"
}

// force string to bytes
func string2byte(str string) []byte  {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}