package main

import (
	"fmt"
	"net"

	"github.com/warent/calzone/cli/v2/cmd"
)

const ADDR = "127.0.0.1"
const PORT = 61895

func main() {
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", ADDR, PORT))
	if err == nil {
		cmd.Execute()
		return
	}
	if err != nil {
		fmt.Println("It looks like Calzone is not running on your machine. Do you want to set it up?")
	}
}
