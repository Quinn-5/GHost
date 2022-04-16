package main

import "fmt"

func Greeting(name string) {
	fmt.Printf("Hello %s, nice to meet you", name)
}

func GetName() string {
	var s string
	fmt.Printf("What is yout name?")
	fmt.Scan(&s)
	return s
}
