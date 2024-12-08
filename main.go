package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/jfyne/wordclouds"
)

func main() {
	topStats, err := os.ReadFile("./top-out")

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	commandLineArgs := os.Args[1:]

	fileContent := strings.Split(string(topStats), "\n")

	contentHeader := strings.Fields(fileContent[6])
	stats := fileContent[7:]

	fmt.Println("Headers:", contentHeader)
	fmt.Println("Sample Stats Line:", strings.Fields(stats[1]))
	fmt.Println("Command Line Args:", commandLineArgs)

	wordCounts := map[string]int{"important": 42, "noteworthy": 30, "meh": 3}

	fontFile, err := os.Open("Roboto-Black.ttf")
	if err != nil {
		log.Fatalf("Error loading font file: %v", err)
	}

	// customize this
	var defaultColors = []color.Color{
		color.RGBA{255, 255, 204, 255},
		color.RGBA{0x48, 0x48, 0x4B, 0xff},
		color.RGBA{0x59, 0x3a, 0xee, 0xff},
		color.RGBA{0x65, 0xCD, 0xFA, 0xff},
		color.RGBA{0x70, 0xD6, 0xBF, 0xff},
	}

	defer fontFile.Close()

	wordCloud, wordCloudError := wordclouds.NewWordcloud(
		wordCounts,
		wordclouds.Font(fontFile),
		wordclouds.Height(2048),
		wordclouds.Width(2048),
		wordclouds.RandomPlacement(true),
		wordclouds.BackgroundColor(color.RGBA{255, 255, 255, 255}),
		wordclouds.Colors(defaultColors),
	)

	if wordCloudError != nil {
		log.Fatalln(wordCloudError)
	}

	image := wordCloud.Draw()

	outputWallpaperFile, err := os.Create("wordcloud.png")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}

	defer outputWallpaperFile.Close()

	if err := png.Encode(outputWallpaperFile, image); err != nil {
		log.Fatalf("Error encoding image to file: %v", err)
	}

	fmt.Println("SAVED")
}
