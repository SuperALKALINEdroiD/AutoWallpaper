package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jfyne/wordclouds"
)

type wordFrequency struct {
	Word      string
	StatValue int
}

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

	var defaultColors = generateColors(50)

	var backgroundColor wordclouds.Option
	if mode == "light" {
		backgroundColor = wordclouds.BackgroundColor(color.RGBA{255, 255, 255, 255})
	} else {
		backgroundColor = wordclouds.BackgroundColor(color.RGBA{0, 0, 0, 255})
	}

	defer fontFile.Close()

	// TODO: dynamic screensize and colors
	wordCloud, wordCloudError := wordclouds.NewWordcloud(
		statsAsMap,
		wordclouds.Font(fontFile),
		wordclouds.FontMaxSize(62),
		wordclouds.FontMinSize(10),
		wordclouds.RandomPlacement(true),
		wordclouds.Height(1200),
		wordclouds.Width(1920),
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

func makeMap(statLines []string, statIndex int) map[string]int {
	wordMap := make(map[string]int)

	for _, statLine := range statLines {
		normalizedContent := strings.Fields(statLine)

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

	var minValue, maxValue int = math.MaxInt, math.MinInt
	for _, value := range wordMap {
		if value < minValue {
			minValue = value
		}
		if value > maxValue {
			maxValue = value
		}
	}

	var frequencies []wordFrequency
	for word, count := range wordMap {
		frequencies = append(frequencies, wordFrequency{Word: word, Count: count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].StatValue < frequencies[j].StatValue
	})

	rankedMap := make(map[string]int)
	for index, item := range frequencies {
		rankedMap[item.Word] = index + 1
	}

	return rankedMap
}

func generateColors(count int) []color.Color {
	colors := make([]color.Color, count)

	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < count; i++ {
		r := uint8(rand.Intn(256))
		g := uint8(rand.Intn(256))
		b := uint8(rand.Intn(256))
		colors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
	}

	return colors
}
