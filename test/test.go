package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	s := "hello"

	s = string(md5.New().Sum([]byte(s)))

	fmt.Println(s)
}
