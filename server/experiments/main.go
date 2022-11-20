package main

import (
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mccutchen/palettor"
)

func RunWithSize(img image.Image, averageSize, iterations int) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	aspectRatio := float64(h) / float64(w)

	resizeWidth := (2 * float64(averageSize)) / (aspectRatio + 1)
	width := int(resizeWidth)
	height := int(resizeWidth * aspectRatio)

	resized := imaging.Resize(img, width, height, imaging.Box)

	start := time.Now()
	p, err := palettor.Extract(3, iterations, resized)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(p.Entries())
	end := time.Now()

	log.Printf(
		"Took %d milliseconds to compute for dimensions (%d, %d)...",
		end.Sub(start).Milliseconds(), width, height,
	)
}

func main() {
	f, err := os.Open("../backgrounds/2.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Testing change in dimensions...")
	RunWithSize(img, 600, 1)
	RunWithSize(img, 800, 1)
	RunWithSize(img, 1080, 1)
	RunWithSize(img, 1440, 1)

	log.Println("Testing change in iterations...")
	RunWithSize(img, 600, 10)
	RunWithSize(img, 600, 20)
	RunWithSize(img, 600, 40)
	RunWithSize(img, 600, 100)
}
