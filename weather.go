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
func getUserInput() string {
	input := bufio.NewReader(os.Stdin)
	s, _ := input.ReadString('\n')
	return s
}

// returns the forecast type, unit type and days requested (if the forecast
// type is future, if not return zero)
func getInput() (string, string, int) {

	// forecast type
	fmt.Print("See current weather or future forecast? [c/f] ")
	forecast_type := getUserInput()
	forecast_type = forecast_type[:len(forecast_type)-1]

	// check days if forecast type is future, then request days
	days := 0
	if forecast_type == "f" {
		fmt.Print("How many days to display? (1 to 7 inclusive) ")
		temp := getUserInput()
		temp = temp[:len(temp)-1]
		days, _ = strconv.Atoi(temp)
	}

	// unit type
	fmt.Print("See information in Imperial or Metric? [i/m] ")
	unit_type := getUserInput()
	unit_type = unit_type[:len(unit_type)-1]

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

// print the current forecast with the data provided
func printCurrentForecast(data map[string]interface{}, unit_type string) {
	// collect information
	current := data["current"].(map[string]interface{})

	// all non-metric data
	last_updated := current["last_updated"].(string)
	condition := current["condition"].(map[string]interface{})["text"].(string)
	cloud := int(current["cloud"].(float64))
	humidity := int(current["humidity"].(float64))

	// metric data
	feels_like := current["feelslike_c"].(float64)
	temperature := current["temp_c"].(float64)
	wind := current["wind_kph"].(float64)
	rainfall := current["precip_mm"].(float64)
	visibility := current["vis_km"].(float64)

	// metric suffixes
	temp_suffix := "C"
	wind_suffix := "Km/h"
	rain_suffix := "mm"
	visibility_suffix := "Km"

	// imperial data
	if unit_type == "i" {
		feels_like = current["feelslike_f"].(float64)
		temperature = current["temp_f"].(float64)
		wind = current["wind_mph"].(float64)
		rainfall = current["precip_in"].(float64)
		visibility = current["vis_miles"].(float64)

		// imperial suffixes
		temp_suffix = "F"
		wind_suffix = "Mi/h"
		rain_suffix = "in"
		visibility_suffix = " Miles"
	}

	// output information
	fmt.Printf("CURRENT WEATHER (last updated %s)\n", last_updated)
	fmt.Printf("%s, %d cloud and it feels like %.1f degrees %s.\n",
		condition, cloud, feels_like, temp_suffix)
	fmt.Printf("Temperature:  %.1fo%s\n", temperature, temp_suffix)
	fmt.Printf("Wind:         %.1f%s\n", wind, wind_suffix)
	fmt.Printf("Rainfall:     %.1f%s\n", rainfall, rain_suffix)
	fmt.Printf("Humidity:     %d%%\n", humidity)
	fmt.Printf("Visibility:   %.1f%s\n", visibility, visibility_suffix)
}

// print the future forecast with the data provided
func printFutureForecast(data map[string]interface{}, unit_type string) {
	// future forecast
	forecast := data["forecast"].(map[string]interface{})
	forecast_days := forecast["forecastday"].([]interface{})
	fmt.Printf("WEATHER FORECAST (next %d days)\n", len(forecast_days))

	for index, value := range forecast_days {
		day := value.(map[string]interface{})
		day_s := day["day"].(map[string]interface{})

		// all non-metric data
		date := day["date"].(string)
		days_away := index + 1
		cond := day_s["condition"].(map[string]interface{})["text"].(string)
		sunrise := day["astro"].(map[string]interface{})["sunrise"].(string)
		sunset := day["astro"].(map[string]interface{})["sunset"].(string)
		humidity := int(day_s["avghumidity"].(float64))

		// metric data
		temp_min := day_s["mintemp_c"].(float64)
		temp_max := day_s["maxtemp_c"].(float64)
		temp_avg := day_s["avgtemp_c"].(float64)
		max_wind := day_s["maxwind_kph"].(float64)
		rainfall := day_s["totalprecip_mm"].(float64)
		visibility := day_s["avgvis_km"].(float64)

		// metric suffixes
		temp_suffix := "C"
		wind_suffix := "Km/h"
		rain_suffix := "mm"
		visibility_suffix := "Km"

		if unit_type == "i" {
			// imperial data
			temp_min = day_s["mintemp_f"].(float64)
			temp_max = day_s["maxtemp_f"].(float64)
			temp_avg = day_s["avgtemp_f"].(float64)
			max_wind = day_s["maxwind_mph"].(float64)
			rainfall = day_s["totalprecip_in"].(float64)
			visibility = day_s["avgvis_miles"].(float64)

			// imperial suffixes
			temp_suffix = "F"
			wind_suffix = "Mi/h"
			rain_suffix = "in"
			visibility_suffix = " Miles"
		}

		fmt.Printf("\n")
		fmt.Printf("|| %s (%d day(s) away) ||\n", date, days_away)
		fmt.Printf("|| %s\n", cond)
		fmt.Printf("|| Sunrise:  %s\n", sunrise)
		fmt.Printf("|| Sunset:   %s\n", sunset)
		fmt.Printf("Temperature\n")
		fmt.Printf("  Max:            %.1fo%s\n", temp_max, temp_suffix)
		fmt.Printf("  Min:            %.1fo%s\n", temp_min, temp_suffix)
		fmt.Printf("  Avg:            %.1fo%s\n", temp_avg, temp_suffix)
		fmt.Printf("Max Wind:         %.1f%s\n", max_wind, wind_suffix)
		fmt.Printf("Total Rainfall:   %.1f%s\n", rainfall, rain_suffix)
		fmt.Printf("Avg Humidity:     %d%%\n", humidity)
		fmt.Printf("Avg Visibility:   %.1f%s\n", visibility, visibility_suffix)
		fmt.Printf("\n")
		fmt.Printf("--------------------------------\n")
	}
}

// use all collected info to generate output (through other functions)
func displayOutput(data map[string]interface{}, forecast_type string,
	unit_type string, forecast_days int) {
	// current forecast
	if forecast_type == "c" {
		printCurrentForecast(data, unit_type)
	} else {
		printFutureForecast(data, unit_type)
	}
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

	// print weather info
	displayOutput(weather_info, forecast_type, unit_type, forecast_days)
}
