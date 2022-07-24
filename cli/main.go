package main

import (
	"fmt"
	"net"
)

const ADDR = "127.0.0.1"
const PORT = 61895

func main() {
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", ADDR, PORT))
	fmt.Println("ERR", err)
	if err == net.ErrClosed {
		fmt.Println("YES")
	}
}
