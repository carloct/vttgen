package main

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
	Start    float64
}

func New() *VttGenerator {
	return &VttGenerator{
		TimeSpan:    10,
		ThumbWidth:  120,
		SpriteWidth: 10,
		Commands: map[string]string{
			"details": "ffprobe -i %s",
			"poster":  "ffmpeg -ss %d -i %s -y -vframes 1 %s/poster.jpg",
			"thumbs":  "ffmpeg -ss %0.04f -i %s -y -an -sn -vsync 0 -q:v 5 -threads 1 -vf scale=%d:-1,select=\"not(mod(n\\,%d))\" \"%s/thumbnails/%s-%%04d.jpg\"",
		},
	}
}

func (v *VttGenerator) Generate(input string, output string, timespan interface{}, thumbwidth interface{}) error {

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return errors.New("Cannot read the input file")
	}
	v.Input = input
	v.Output = output

	if timespan != nil {
		v.TimeSpan = timespan.(int)
	}

	if thumbwidth != nil {
		v.ThumbWidth = thumbwidth.(int)
	}

	_ = v.Details()

	_ = v.Poster(20)

	return nil

}

func (v *VttGenerator) Details() error {
	out, err := exec.Command("ffprobe", "-i", v.Input).CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return err
	}

	v.Video.Duration, v.Video.Start = duration(string(out))
	v.Video.Tbr = tbr(string(out))

	return nil
}

func (v *VttGenerator) Poster(frame int) error {
	alpha := strconv.Itoa(frame)
	_, err := exec.Command("ffmpeg", "-ss", alpha, "-i", v.Input, "-y", "-vframes", "1", v.Output+"poster.jpg").CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func duration(output string) (int, float64) {
	r, _ := regexp.Compile(`(?is)Duration: ((\d+):(\d+):(\d+))\.\d+, start: ([^,]*)`)
	matches := r.FindStringSubmatch(string(output))
	start, _ := strconv.ParseFloat(matches[5], 64)
	return timeToSeconds(matches[1]), start
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
