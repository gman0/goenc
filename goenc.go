package main

import (
	"fmt"
	"github.com/gman0/goenc/enc"
)

func main() {
	keys := enc.GenerateKeyPair()
	fmt.Println(string(keys.Private))
	fmt.Println(string(keys.Public))
}
