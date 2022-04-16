package main

import "fmt"

func Greeting(name string) {
	fmt.Printf("Hello %s, nice to meet you\n", name)
}

func GetName() string {
	var s string
	fmt.Println("What is your name?")
	fmt.Scan(&s)
	return s
}
