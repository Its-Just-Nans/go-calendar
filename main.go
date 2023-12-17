// create basic main go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func prettyPrint(data map[string]interface{}) {
	customPrintln("{")
	for key, value := range data {
		switch v := value.(type) {
		case float64:
			customPrintf(`    "%s": "%.2f",%s`, key, v, "\n")
		default:
			customPrintf(`    "%s": "%v",%s`, key, v, "\n")
		}
	}
	customPrintln("}")
}

func mod(a int, b int) int {
	return (a%b + b) % b
}

func parsedData(result []map[string]interface{}, keyName string, counterKey string) (int, map[time.Time]int, time.Time, time.Time) {
	max := 1
	dateData := make(map[time.Time]int)
	var lowest time.Time
	var highest time.Time
	for _, oneData := range result {
		currentDate, ok := oneData[keyName].(string)
		if !ok {
			customPrintf(`Error: Key "%s" not found in:%s`, keyName, "\n")
			prettyPrint(oneData)
			continue
		}
		layout := time.RFC3339[:len(currentDate)]
		parsedDate, err := time.Parse(layout, currentDate)
		if err != nil {
			customPrintf(`Error: during date parsing of "%s" in:%s`, currentDate, "\n")
			prettyPrint(oneData)
			customPrintln(err)
			continue
		}
		if highest.IsZero() || lowest.IsZero() {
			lowest = parsedDate
			highest = parsedDate
		} else if parsedDate.After(highest) {
			highest = parsedDate
		} else if parsedDate.Before(lowest) {
			lowest = parsedDate
		}
		if counterKey == "" {
			if dateData[parsedDate] == 0 {
				dateData[parsedDate] = 1
			} else {
				dateData[parsedDate]++
				if dateData[parsedDate] > max {
					max = dateData[parsedDate]
				}
			}
		} else {
			value, ok := oneData[counterKey].(float64)
			if !ok {
				customPrintf(`Error: No key "%s" (counter key) in: %s`, counterKey, "\n")
				prettyPrint(oneData)
				continue
			}
			if dateData[parsedDate] == 0 {
				dateData[parsedDate] = int(value)
			} else {
				dateData[parsedDate] += int(value)
			}
			if dateData[parsedDate] > max {
				max = dateData[parsedDate]
			}
		}
	}
	return max, dateData, lowest, highest
}

func AddGroupText(counter int, currentDay time.Time, cubeSize int, marginSize int, str *strings.Builder) {
	yTranslate := cubeSize * 2
	yTranslate += int(float32(counter) * float32(cubeSize*7+marginSize*7) * 1.5)
	xTranslate := 50
	str.WriteString(fmt.Sprintf(`<text font-size="%d" font-family="system-ui, monospace" x="%d" y="%d">%d</text>%s`, cubeSize, 26*cubeSize+26*marginSize+xTranslate, yTranslate, currentDay.Year(), "\n"))
	date := time.Time{}
	date = date.AddDate(0, 0, int(firstDayOfWeek)-1)
	str.WriteString(fmt.Sprintf("<g transform=\"translate(%d,%d)\">\n", 10, yTranslate+cubeSize*2))
	for i := 0; i < 4; i++ {
		t := fmt.Sprintf(`<text font-size="%d" font-family="system-ui, monospace" x="%d" y="%d">%s</text>%s`, cubeSize, 0, (i*cubeSize+i*marginSize)*2, date.Format("Mon"), "\n")
		str.WriteString(t)
		date = date.AddDate(0, 0, 2)
	}
	str.WriteString("</g>\n")
	str.WriteString(fmt.Sprintf("<g transform=\"translate(%d,%d)\">\n", xTranslate, yTranslate))
}

func AddStyle(str *strings.Builder, max int) {
	str.WriteString("<style>")
	for i := 0; i <= max; i++ {
		color := generateDarkenedColor(i, max)
		s := fmt.Sprintf(`[data-level='%d'] {fill: #%02x%02x%02x;}`, i, color.R, color.G, color.B)
		str.WriteString(s)
	}
	str.WriteString("</style>\n")
}

func generateDarkenedColor(level int, maxLevel int) color.RGBA {
	// Calculate the darkening factor based on the level
	if level == 0 {
		return baseColor
	}
	darkeningFactor := float64(level+1) / float64(maxLevel)
	darkeningFactorR := 0.0
	darkeningFactorG := 0.0
	darkeningFactorB := 0.0
	if baseColor.R > 0 {
		darkeningFactorR = darkeningFactor
	}
	if baseColor.G > 0 {
		darkeningFactorG = darkeningFactor
	}
	if baseColor.B > 0 {
		darkeningFactorB = darkeningFactor
	}
	darkenedR := uint8(float64(baseColor.R) * (1 - darkeningFactorR))
	darkenedG := uint8(float64(baseColor.G) * (1 - darkeningFactorG))
	darkenedB := uint8(float64(baseColor.B) * (1 - darkeningFactorB))

	return color.RGBA{darkenedR, darkenedG, darkenedB, baseColor.A}
}

var cubeSize = 15
var marginSize = 2
var baseColor = color.RGBA{0, 200, 200, 255}

var firstDayOfWeek = time.Monday
var keyName = "date"
var counterKey = ""
var filename = "data.json"
var outputFilename = "out.svg"
var quiet = false

func customPrintln(a ...any) {
	if !quiet {
		fmt.Println(a...)
	}
}

func customPrintf(format string, a ...any) {
	if !quiet {
		fmt.Printf(format, a...)
	}
}

func parseFlag() {
	day := int(firstDayOfWeek)
	flag.StringVar(&filename, "i", filename, "Input JSON file")
	flag.StringVar(&keyName, "k", keyName, "Key name for date")
	flag.IntVar(&day, "d", day, "First day of the week")
	flag.StringVar(&counterKey, "c", counterKey, "Key name for counter")
	flag.BoolVar(&quiet, "q", quiet, "Quiet mode")
	flag.StringVar(&outputFilename, "o", outputFilename, "Output filename")
	flag.Parse()
	day = mod(day, 7)
	firstDayOfWeek = time.Weekday(day)
}

func readData() []byte {
	jsonFile, err := os.Open(filename)
	if err != nil {
		customPrintln(err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		customPrintln(err)
		os.Exit(1)
	}
	return byteValue
}

func daysInFirstWeek(year int) int {
	firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	i := 0
	for firstDayOfYear.Weekday() != firstDayOfWeek {
		firstDayOfYear = firstDayOfYear.AddDate(0, 0, 1)
		i++
	}
	return i
}

func generateSvg(result *[]map[string]interface{}) []byte {
	max, dateData, lowest, highest := parsedData(*result, keyName, counterKey)
	currentDay := time.Date(lowest.Year(), 01, 01, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(highest.Year(), 12, 31, 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
	var str strings.Builder
	AddStyle(&str, max)
	defaultStr := fmt.Sprintf(`<rect data-date="%%s" width="%d" height="%d" x="%%d" y="%%d" data-level="%%d">`, cubeSize, cubeSize)
	counter := 0
	AddGroupText(counter, currentDay, cubeSize, marginSize, &str)
	counterWeek := 0
	nbDaysInFirstWeek := 0
	for currentDay != lastDay {
		dayOfWeek := int(currentDay.Weekday())
		formattedDate := currentDay.Format("2006-01-02")
		if currentDay.YearDay() == 1 {
			counter++
			nbDaysInFirstWeek = daysInFirstWeek(currentDay.Year())
			counterWeek = 0
		}
		x := counterWeek*cubeSize + counterWeek*marginSize
		if dayOfWeek == mod(int(firstDayOfWeek)-1, 7) && currentDay.YearDay() > nbDaysInFirstWeek-1 {
			counterWeek++
		}
		dayOfWeek = mod(dayOfWeek-int(firstDayOfWeek), 7) + 1
		y := dayOfWeek*cubeSize + dayOfWeek*marginSize
		toDisplay := 0
		if val, ok := dateData[currentDay]; ok {
			toDisplay = val
		}
		toAdd := fmt.Sprintf(defaultStr, formattedDate, x, y, toDisplay)
		title := fmt.Sprintf(`<title>%s - %s: %d</title>`, currentDay.Format("Mon"), formattedDate, toDisplay)
		final := toAdd + "\n" + title + "</rect>\n"
		str.WriteString(final)
		previousY := currentDay.Year()
		currentDay = currentDay.AddDate(0, 0, 1)
		if currentDay.Year() != previousY {
			str.WriteString("</g>\n")
			if currentDay != lastDay {
				AddGroupText(counter, currentDay, cubeSize, marginSize, &str)
			}
		}
	}
	height := counter * (cubeSize*10 + marginSize*18)
	width := 53*cubeSize + 53*marginSize + 75
	header := fmt.Sprintf(`<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">%s`, width, height, width, height, "\n")
	str.WriteString("</svg>\n")
	data := []byte(header + str.String())
	return data
}

func parseData(byteValue *[]byte) []map[string]interface{} {
	var result []map[string]interface{}
	err := json.Unmarshal(*byteValue, &result)
	if err != nil {
		customPrintln("Error during JSON parsing")
		customPrintln(err)
		os.Exit(1)
	}
	if len(result) == 0 {
		customPrintf("JSON array of '%s' is empty\n", filename)
		os.Exit(1)
	}
	return result
}

func main() {
	parseFlag()
	byteValue := readData()
	result := parseData(&byteValue)
	data := generateSvg(&result)
	err := os.WriteFile(outputFilename, data, 0644)
	if err != nil {
		customPrintln(err)
		os.Exit(1)
	}
}
