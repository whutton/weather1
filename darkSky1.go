
//    program name darkSky.go
//    This is set for the whutton@gmail.com key and the locations for Woodstock, GA
//
//    Bill Hutton
//    7-26-18
//
package main

import (
	"net/http"
	"fmt"
	"time"
	"io/ioutil"
	"os"
	"encoding/json"
	//"strconv"
	"log"
)

var programName = "darkSky1"
var fullTimeLayout = "Mon Jan 2 15:04:05 -0700 MST 2006"
var hourTimeLayout = "Mon Jan 2 15:04:05 -0700 MST 2006"


type DarkResponse struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitudelatitude"`
	Timezone string `json:"timezone"`
	Currently CurValues `json:"currently"`
	Minutely MinValues `json:"minutely"`
	Hourly   HourValues `json:"hourly"`
	Daily   []DailyValues `json:"daily"` 
	Alert  []AlertValues `json:"alerts"` 
}

type CurValues struct {
	Time int64 `json:"time"`
	Summary string `json:"summary"`
	Icon string `json:"icon"`
	NearestStormDistance float64 `json:"nearestStormDistance"`
	PrecipIntensity float64	`json:"precipIntensity"`
	PrecipIntensityError float64	`json:"precipIntensityError"`
	PrecipProbability float64	`json:"precipProbability"`
	PrecipType  string	`json:"precipType"`
	Temperature float64	`json:"temperature"`
	ApparentTemperature float64	`json:"apparentTemperature"`
	DewPoint float64	`json:"dewPoint"`
	Humidity float64	`json:"humidity"`
	Pressure float64	`json:"pressure"`
	WindSpeed  float64	`json:"windSpeed"`
	WindGust float64	`json:"windGust"`
	WindBearing float64	`json:"windBearing"`
	CloudCover float64	`json:"cloudCover"`
	UvIndex float64	 `json:"uvIndex"`
	visibility float64	`json:"visibility"`
	Ozone float64	`json:"ozone"`
}

type SpecificValues struct {
	Time	int64	`json:"time"`
	Summary	string	`json:"summary"`
	Icon	string	`json:"icon"`
	PrecipIntensity	float64	`json:"precipIntensity"`
	PrecipProbability	float64	`json:"precipProbability"`
	PrecipType	string	`json:"precipType"`
	Temperature	float64	`json:"temperature"`
	ApparentTemperature	float64	`json:"apparentTemperature"`
	DewPoint	float64	`json:"dewPoint"`
	Humidity	float64	`json:"humidity"`
	Pressure	float64	`json:"pressure"`
	WindSpeed	float64	`json:"windSpeed"`
	WindGust	float64	`json:"windGust"`
	WindBearing	float64	`json:"windBearing"`
	CloudCover	float64	`json:"cloudCover"`
	UvIndex	float64	`json:"uvIndex"`
	Visibility	float64	`json:"visibility"`
	Ozone	float64	`json:"ozone"`
	
	TemperatureHigh	float64	`json:"temperatureHigh"`
	TemperatureHighTime	int64	`json:"temperatureHighTime"`
	TemperatureLow	float64 	`json:"temperatureLow"`
	TemperatureLowTime	int64	`json:"temperatureLowTime"`

}

type MinValues struct {
	Summary string `json:"summary"`
	Icon string `json:"icon"`
	Details []SpecificValues `json:"data"`
	

}

type HourValues struct {
	Summary string `json:"summary"`
	Icon string `json:"icon"`
	Details []SpecificValues `json:"data"`

}

type DailyValues struct {
	Summary string `json:"summary"`
	Icon string `json:"icon"`
	Details []SpecificValues `json:"data"`

}

type AlertValues struct {
	Title string `json:"title"`
	Time int64 `json:"time"`
	Expires int64 `json:"expires"`   // time
	Description string `json:"description"`
	Uri string `json:"URI"`

}

func main () {


	//
	arg := os.Args
	listenPort := arg[1];
	string1 := fmt.Sprintf("listen for weather request on %s",listenPort)
	fmt.Println(string1)
	
	http.HandleFunc("/", handlerIn)
	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}


func handlerIn(w http.ResponseWriter, r *http.Request) {
	retString := runDark()
	fmt.Fprintf(w,"%s",retString)
}

func runDark() string {
	outString := ""
	resp, err := http.Get("https://api.darksky.net/forecast/oscar/34.126681,-84.55377")
	checkError(err)
	var st1 = new(DarkResponse)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &st1)

	outString += "<!DOCTYPE html>"
	outString += "<html><head><meta = charset = \"utf-8\">"
	outString += "<title>Darksky</title>"
	outString += "</head>"
	outString += "<body>"


	queryCurrent := st1.Currently
	tour1 := returnTime(queryCurrent.Time,"full")
	outString += "<H2>"+tour1+"</H2>"
	str1 := fmt.Sprintf("%2.3f",queryCurrent.Temperature)
	outString += "<H3>  Temp: "+str1+"   "+queryCurrent.Summary+"</H3>"

	
	outString += "</tr></table>"
	
	hourOne := st1.Hourly
	hourDetails := hourOne.Details
	outString += "<H4>"+"This hour " + hourOne.Summary +  "</H4>"
	
	outString += "<table style=\"width:45%\">"
	outString += "<tr>"
	outString += "<th>Date</th><th>Temp</th><th>Summary</th><th>Intensity</th><th>Probability</th>"
	outString += "<th>Humidity</th><th>Windspeed</th><th>Cloudcover</th>"
	outString += "</tr>"

	hourLength := len(hourDetails)
	holdDay := -1
	lowTemp := 0.00
	highTemp := 0.00
	lowHigh := ""
	for i := 0; i < hourLength; i++ {
		ele1 := hourDetails[i]
		dateElement := returnTime(ele1.Time,"short")
		dateHold := dateElement
		temp1 := time.Unix(ele1.Time,0)
		dayNumber := temp1.Day() 

		if i == 0 {

			lowTemp, highTemp = getLowHigh(0,hourDetails,dayNumber,hourLength)
			lowHigh = fmt.Sprintf("l) %2.1f h) %2.1f",lowTemp,highTemp)
			//fmt.Println(lowHigh)
		}

		if dayNumber != holdDay {
			str2 := fmt.Sprintf("%d",dayNumber)
			if i == 0 {
				outString += "<tr><td>Day: "+str2+"</td><td></td><td>"+lowHigh+"</td></tr>"
			}
			lowHigh = ""
			holdDay = dayNumber
			if i != 0 {

				lowTemp, highTemp = getLowHigh(i,hourDetails,holdDay,hourLength)
				lowHigh = fmt.Sprintf("l) %2.1f h) %2.1f",lowTemp,highTemp)				
				outString += "<tr><td>Day: "+str2+"</td><td></td><td>"+lowHigh+"</td></tr>"

			}
		}

		intensity1 := ele1.PrecipIntensity
		prob1 := 100 * ele1.PrecipProbability

		humidity := ele1.Humidity * 100
		windspeed := ele1.WindSpeed
		cloudcover := ele1.CloudCover * 100
		summary1 := ele1.Summary

		outString += "<tr>"
		outString += "<td>"+dateHold+"</td>"

		str1 := fmt.Sprintf("%2.3f",ele1.Temperature)

		
		outString += "<td" + ">"+str1+"</td>"
		
		outString += "<td>"+summary1+"</td>"

		str1 = fmt.Sprintf("%3.0f",intensity1)
		outString += "<td>"+str1+"</td>"

		str1  = fmt.Sprintf("%.3f",prob1)
		outString += "<td>"+str1+"</td>"

		str1 = fmt.Sprintf("%3.0f",humidity)
		outString += "<td>"+str1+"</td>"
		str1 = fmt.Sprintf("%2.3f",windspeed)
		outString += "<td>"+str1+"</td>"
		str1 = fmt.Sprintf("%2.3f",cloudcover)
		outString += "<td>"+str1+"</td>"
		outString +=("</tr>")

	}
	outString += "</tr></table>"
	outString += "<footer>Bill Hutton  8-2-18</footer>"
	outString += "</body></html>"
	return outString
}


func checkError(err error) {
	if err != nil {
		fmt.Println(programName + ": Fatal error ", err.Error())
	}
}

func returnTime(inValue int64,inType string) string {
	tm := time.Unix(inValue, 0)
	if inType == "full" {
		tReturn := tm.Format(fullTimeLayout)
		return tReturn
	}
	tReturn := tm.Format("3:04PM")
	return tReturn
}

func getLowHigh(inStart int,hours []SpecificValues,holdDay int,totalDays int) (float64, float64) {	
	holdLow := 99999.99
	holdHigh := -99999.99

	for i := inStart; i < totalDays;i++ {
		ele1 := hours[i]

		temp1 := ele1.Temperature
		if temp1 < holdLow {
			holdLow = temp1
		}
		if temp1 > holdHigh {
			holdHigh = temp1			
		}
		temp2 := time.Unix(ele1.Time,0)
		dayNumber := temp2.Day()
		if dayNumber != holdDay {
			return holdLow, holdHigh
		}
	}
        return holdLow,holdHigh
}
