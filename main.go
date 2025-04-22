package main

import (
	"fmt"

	bk "github.com/maxime-peim/gokite/pkg"
)

func main() {
	socketPath := "bird.ctl"
	kite, err := bk.NewBirdKite(socketPath)
	if err != nil {
		panic(err)
	}

	status, err := kite.Status()
	if err != nil {
		panic(err)
	}
	fmt.Println("Status:", status)
}
