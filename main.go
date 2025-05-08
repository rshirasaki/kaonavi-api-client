package main

import (
	"fmt"
	"os"
)

func main() {
	client := NewClient(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	fmt.Println(client.GetMembers())
}
