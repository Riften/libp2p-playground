package main

import (
	"fmt"
	"main/cmd"
)

func main() {

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

}
