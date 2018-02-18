package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// uses provided weather api to receive information and return it
func useWeatherApi(url string) map[string]interface{} {
	var data map[string]interface{} // will hold the data

	client := http.Client{}
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	res, _ := client.Do(request)
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &data)

	return data
}

// returns the weather information based on the input, after creating url
func getWeatherInfo(forecast_type string, unit_type string, forecast_days int,
	api_key string) map[string]interface{} {

	url := "https://extension.codecadets.com/api/current/?key=" + api_key
	if forecast_type == "f" {
		url = "https://extension.codecadets.com/api/forecast/?key=" + api_key
		url += "&days=" + strconv.Itoa(forecast_days)
	}

	response := useWeatherApi(url)
	return response
}

func main() {
	api_key := "fb171458a2e430573cd3edc90962b473"

	// introduction
	fmt.Println("Welcome to ASCIIWeather!")

	// get input
	forecast_type, unit_type, forecast_days := getInput()

	// create api url and retrieve info
	weather_info := getWeatherInfo(forecast_type, unit_type, forecast_days,
		api_key)

	fmt.Println(weather_info)
}
