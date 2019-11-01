package main

import (
	"fmt"
	"log"
)

func testFatalError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		log.Fatal(e.Error())
	}
}

func testQueryError(e error, message string) {
	if e != nil {
		fmt.Println(message, e.Error())
	}
}
