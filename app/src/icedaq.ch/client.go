package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	var apiBaseUrl string

	value := os.Getenv("APIURL")
	if len(value) == 0 {
		apiBaseUrl = "http://localhost:9001/api"
	} else {
		apiBaseUrl = value
	}

	fmt.Println("Using api url: " + apiBaseUrl)

	for {
		start := time.Now()

		// Get all artist.
		resp, err := http.Get(apiBaseUrl + "/artists")
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		var f interface{}
		err = json.Unmarshal(contents, &f)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		entries := f.([]interface{})
		var maxId float64
		maxId = 0
		for _, entry := range entries {
			mymap := entry.(map[string]interface{})
			currentId := mymap["id"].(float64)
			if currentId >= maxId {
				maxId = currentId
			}
		}
		//fmt.Println(maxId)

		// Pick one at random.
		artistId := random(1, int(maxId))

		// Get shows. Pick one at random.
		resp, err = http.Get(apiBaseUrl + "/shows/" + strconv.Itoa(artistId))
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		contents, err = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(contents, &f)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		entries = f.([]interface{})
		var showIds []int
		for _, entry := range entries {
			mymap := entry.(map[string]interface{})
			showIds = append(showIds, int(mymap["id"].(float64)))
		}
		//fmt.Println(maxId)

		// Pick one at random.
		showIdPos := random(0, len(showIds))
		showId := showIds[showIdPos]

		// Get ticket info. We will not actually do anything with this data but just to simulate the client.
		resp, err = http.Get(apiBaseUrl + "/tickets/" + strconv.Itoa(showId))
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		contents, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		// Buy ticket.
		resp, err = http.Get(apiBaseUrl + "/buy/" + strconv.Itoa(showId))
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		contents, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		// See if the purchase was successfull!
		//fmt.Println(string(contents))
		if string(contents) == "true" {
			fmt.Println("Bought ticket for show " + strconv.Itoa(showId))
		} else {
			fmt.Println("Purchase failed for show " + strconv.Itoa(showId))
		}

		elapsed := time.Since(start)
		fmt.Printf("Purchase took %s\n", elapsed)

		// Log this to stackdriver.
	}

	//os.Exit(0)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
