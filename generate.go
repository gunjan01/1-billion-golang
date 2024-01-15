package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

var europeanCities = []string{
	"Paris", "London", "Berlin", "Madrid", "Rome", "Amsterdam", "Vienna", "Athens", "Prague", "Warsaw",
	"Lisbon", "Budapest", "Stockholm", "Oslo", "Copenhagen", "Dublin", "Brussels", "Helsinki", "Ljubljana", "Zagreb",
	"Bucharest", "Sofia", "Bratislava", "Luxembourg", "Dubrovnik", "Barcelona", "Munich", "Milan", "Geneva", "Edinburgh",
	"Reykjavik", "Tallinn", "Vilnius", "Bruges", "Porto", "Valletta", "Zurich", "Krakow", "Nice", "Cologne", "Lyon",
	"Frankfurt", "Hamburg", "Malaga", "Naples", "Salzburg", "Belfast", "Bilbao", "Glasgow", "Seville", "Marseille",
	"Gothenburg", "Minsk", "Belgrade", "Sarajevo", "Skopje", "Tirana", "Riga", "Vaduz", "Andorra la Vella", "Monaco",
	"Nicosia", "San Marino", "Vatican City", "Podgorica", "Bern", "Amman", "Yerevan", "Baku", "Chisinau", "Tbilisi",
	"Kiev", "Minsk", "Tirana", "Vilnius", "Bucharest", "Sofia", "Bratislava", "Ljubljana", "Zagreb", "Podgorica", "Sarajevo",
}

func generateRandomCity() string {
	rand.Seed(time.Now().UnixNano())
	return europeanCities[rand.Intn(len(europeanCities))]
}

func generateRandomTemperature() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*20 - 5 // Generates temperatures between -5 and 15 degrees Celsius
}

func Generate() {
	file, err := os.OpenFile("measurements.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1000000000; i++ {
		city := generateRandomCity()
		temp := strconv.Itoa(int(generateRandomTemperature()))

		s := city + ";" + temp

		_, err = file.WriteString(s + "\n")
		if err != nil {
			panic(err)
		}
	}

	file.Close()

}
