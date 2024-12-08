package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jfyne/wordclouds"
)

func main() {
	topStats, err := os.ReadFile("./top-out")

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var mode string
	var statIndex int

	commandLineArgs := os.Args[1:]
	if len(commandLineArgs) > 0 {
		// assign mode, color and stat type
		for _, arg := range commandLineArgs {
			argOption := strings.Split(arg, "=")

			switch argOption[0] {
			case "type":
				if argOption[1] == "cpu" {
					statIndex = 8 // cpu
				} else {
					statIndex = 9 // memory
				}
			case "color":
				fmt.Println("two")
			case "mode":
				mode = argOption[1]
			}
		}
	}

	fileContent := strings.Split(string(topStats), "\n")

	contentHeader := strings.Fields(fileContent[6])
	stats := fileContent[7:]

	fmt.Println("Headers:", contentHeader)
	fmt.Println("Sample Stats Line:", strings.Fields(stats[1]))

	statsAsMap := makeMap(stats, statIndex)

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

	var backgroundColor wordclouds.Option

	if mode == "light" {
		backgroundColor = wordclouds.BackgroundColor(color.RGBA{255, 255, 255, 255})
	} else {
		backgroundColor = wordclouds.BackgroundColor(color.RGBA{0, 0, 0, 255})
	}

	defer fontFile.Close()

	wordCloud, wordCloudError := wordclouds.NewWordcloud(
		statsAsMap,
		wordclouds.Font(fontFile),
		wordclouds.RandomPlacement(true),
		wordclouds.Height(2048),
		wordclouds.Width(2048),
		wordclouds.Colors(defaultColors),
		backgroundColor,
	)

	if wordCloudError != nil {
		log.Fatalln(wordCloudError)
	}

	image := wordCloud.Draw()

	outputWallpaperFile, err := os.Create("wallpaper.png")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}

	defer outputWallpaperFile.Close()

	if err := png.Encode(outputWallpaperFile, image); err != nil {
		log.Fatalf("Error encoding image to file: %v", err)
	}

	fmt.Println("SAVED")
}

func makeMap(statLines []string, statIndex int) (wordMap map[string]int) {
	wordMap = make(map[string]int)

	for _, statLine := range statLines {
		normalizedContent := strings.Fields(statLine)

		fmt.Println(normalizedContent)

		if len(normalizedContent) == 0 {
			continue
		}

		statValue, conversionError := strconv.ParseFloat(normalizedContent[statIndex], 64)
		if conversionError != nil {
			fmt.Println(conversionError)
			continue
		}

		statName := normalizedContent[len(normalizedContent)-1]
		wordMap[statName] += int(statValue)
	}

	return wordMap
}
