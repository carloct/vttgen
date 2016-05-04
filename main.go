package main

import "log"

func main() {
	err := Generate("sample.mp4", "./test/test.mp4")
	if err != nil {
		log.Fatal(err)
	}
}
