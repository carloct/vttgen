package vttgen

import "log"

func main() {
	v := New()
	err := v.Generate("sample.mp4", "./test/test.mp4")
	if err != nil {
		log.Fatal(err)
	}
}
