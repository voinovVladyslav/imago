package main

import (
	"fmt"
	"imago/internal/storage"
)


func main() {
	err := storage.Run()
	if err != nil {
		fmt.Println("got error:", err.Error())
	}
}
