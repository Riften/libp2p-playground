package main

import (
	"fmt"
	"github.com/Riften/libp2p-playground/cmd"
)

func main() {

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

}
