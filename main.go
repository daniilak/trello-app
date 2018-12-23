package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Config structure
type Config struct {
	ClientKey string
	AppKey    string
	BoardID   string
}

// Card structure
type Card struct {
	ID   string
	Name string
}

func main() {
	getListCards()
}

func getIDCard() {
	fmt.Print("Enter a number card:")
	var input float64
	fmt.Scanf("%f", &input)
}

func readBody(srcURL string) ([]byte, error) {
	resp, err := http.Get(srcURL)
	if err != nil {
		return nil, fmt.Errorf("with url: %s: %v", srcURL, err)
	}
	defer closeBody(resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("with read body: %s: %v", resp.Body, err)
	}

	return body, nil
}

func closeBody(resp *http.Response) {
	if errClose := resp.Body.Close(); errClose != nil {
		log.Println(errClose)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getListCards() {

	dat, err := ioutil.ReadFile("config.json")
	check(err)
	var config Config
	err = json.Unmarshal(dat, &config)

	var request = strings.Join(
		[]string{
			"https://api.trello.com/1/boards/",
			string(config.BoardID),
			"/cards/?limit=10&fields=name&key=",
			string(config.AppKey),
			"&token=",
			string(config.ClientKey)},
		"")
	respBytes, err := readBody(request)

	var cards []Card
	err = json.Unmarshal(respBytes, &cards)
	if err != nil {
		fmt.Printf("with creating struct from bytes : %s", err)
	}
	for i := 0; i < len(cards); i++ {
		fmt.Printf("%d %+v\n", i+1, cards[i].Name)
	}

}
