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
	Video       Video
	Input       string
	Output      string
	TimeSpan    int
	ThumbWidth  int
	SpriteWidth int
	Commands    map[string]string
}

type Video struct {
	Duration int
	Tbr      int
}

func New() *VttGenerator {
	return &VttGenerator{
		Commands: map[string]string{
			"details": "ffprobe -i %s 2>&1",
			"poster":  "ffmpeg -ss %d -i %s -y -vframes 1 %s/poster.jpg 2>&1",
			"thumbs":  "ffmpeg ss %0.04f -i %s -y -an -sn -vsync 0 -q:v 5 -threads 1 -vf scale=%d:-1,select=\"not(mod(n\\,%d))\" \"%s/thumbnails/%s-%%04d.jpg\" 2>&1",
		},
	}
}

func (v *VttGenerator) Generate(input string, output string) error {

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return errors.New("Cannot read the input file")
	}

	out, err := exec.Command("ffprobe", "-i", input).CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return err
	}

	v.Video.Duration = duration(string(out))
	v.Video.Tbr = tbr(string(out))

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
