package main

import (
	"fmt"
	shared "shared/pkg"
)

func main() {
	shared.InitSockets()
	fmt.Println(len(shared.Sockets))
}
