package main

import "testing"

func TestServer(t *testing.T) {
	server := New(&Config{
		Host: "localhost",
		Port: "3333",
	})
	server.Run()
}