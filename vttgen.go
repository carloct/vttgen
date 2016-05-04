package vttgen

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type VttGenerator struct {
	Duration int
	Tbr      int
}

func New() *VttGenerator {
	return &VttGenerator{}
}

func (v *VttGenerator) Generate(input string, output string) error {

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return errors.New("Cannot read the input file")
	}

	out, err := exec.Command("ffprobe", "-i", input).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	v.Duration = duration(string(out))
	v.Tbr = tbr(string(out))

	return nil

}

func duration(output string) int {
	r, _ := regexp.Compile(`(?is)Duration: ((\d+):(\d+):(\d+))\.\d+, start: ([^,]*)`)
	matches := r.FindStringSubmatch(string(output))
	return timeToSeconds(matches[1])
}

func timeToSeconds(time string) int {
	r, _ := regexp.Compile(`((\d+):(\d+):(\d+))`)
	matches := r.FindStringSubmatch(string(time))

	hours, _ := strconv.Atoi(matches[2])
	total := hours * 3600

	minutes, _ := strconv.Atoi(matches[3])
	total += minutes * 60

	seconds, _ := strconv.Atoi(matches[4])
	total += seconds

	return total
}

func tbr(output string) int {
	r, _ := regexp.Compile(`\b(\d+(?:\.\d+)?) tbr\b`)
	matches := r.FindStringSubmatch(string(output))

	i, _ := strconv.Atoi(matches[1])

	return i
}
