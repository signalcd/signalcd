package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("checking the status...")

	i := 0
	for {
		fmt.Println(i)
		i++

		if i < 60 {
			time.Sleep(time.Second)
		} else {
			return
		}
	}
}
