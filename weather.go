package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// removes trailing newline of a string
func removeTrailingNewline(s string) string {
	return s[:len(s)-1]
}

// returns the forecast type, unit type and days requested (if the forecast
// type is future, if not return zero)
func getInput() (string, string, int) {

	// forecast type
	fmt.Print("See current weather or future forecast? [c/f] ")
	input := bufio.NewReader(os.Stdin)
	forecast_type, _ := input.ReadString('\n')
	forecast_type = removeTrailingNewline(forecast_type)

	// check days if forecast type is future, then request days
	days := 0
	if forecast_type == "f" {
		fmt.Print("How many days to display? (1 to 10 inclusive) ")
		input := bufio.NewReader(os.Stdin)
		temp, _ := input.ReadString('\n')
		temp = removeTrailingNewline(temp)
		days, _ = strconv.Atoi(temp)
	}

	// unit type
	fmt.Print("See information in Imperial or Metric? [i/m] ")
	reader := bufio.NewReader(os.Stdin)
	unit_type, _ := reader.ReadString('\n')
	unit_type = removeTrailingNewline(unit_type)

	return forecast_type, unit_type, days
}

func main() {
	fmt.Println("Welcome to ASCIIWeather!")
	forecast_type, unit_type, forecast_days := getInput()
}
