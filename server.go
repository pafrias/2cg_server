package main

import (
	"fmt"
	"log"
)

func (s *Server) testFatalError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		log.Fatal(e.Error())
	}
}

func (s *Server) testQueryError(e error, message string) {
	if e != nil {
		fmt.Println(message, e.Error())
	}
}
