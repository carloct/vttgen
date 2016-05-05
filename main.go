package main

import "log"

func main() {
	v := New()
	err := v.Generate("sample.mp4", "./test/test.mp4", 10, 23)
	if err != nil {
		log.Fatal(err)
	}
}
