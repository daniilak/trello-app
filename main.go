package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// Member structure
type Member struct {
	ID       string
	IDMember string
}

// Board structure
type Board struct {
	ID             string
	URL            string
	Name           string
	Desc           string
	DescData       string
	Closed         bool
	IDOrganization string
}

func getConfigList() (config Config) {
	dat, err := ioutil.ReadFile("config.json")
	check(err)
	err = json.Unmarshal(dat, &config)
	return config
}

var config = getConfigList()

func main() {
	id := getMemberID()
	args := os.Args[1:]
	if len(args) > 0 {
		switch par := args[0]; par {
		case "-b":
			getListBoards(id)
			break
		case "-c":
			if len(args[1]) > 0 {
				getListCards(args[1])
			} else {
				getListCards("1")
			}
			break
		}
	}
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

func getListCards(id string) {
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

func getListBoards(IDMember string) {
	var request = strings.Join(
		[]string{
			"https://api.trello.com/1/members/",
			IDMember,
			"/boards?filter=all&fields=all&lists=none&memberships=none&organization=false&organization_fields=name%2CdisplayName&key=",
			string(config.AppKey),
			"&token=",
			string(config.ClientKey)},
		"")
	respBytes, err := readBody(request)

	if err != nil {
		fmt.Printf("err: %s", err)
	}
	var boards []Board
	err = json.Unmarshal(respBytes, &boards)
	for i := 0; i < len(boards); i++ {
		fmt.Printf("%d %+v\n", i+1, boards[i].Name)
	}
}

func getMemberID() (id string) {
	var request = strings.Join(
		[]string{
			"https://api.trello.com/1/tokens/",
			string(config.ClientKey),
			"?token=",
			string(config.ClientKey),
			"&key=",
			string(config.AppKey)},
		"")
	respBytes, err := readBody(request)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	var member Member
	err = json.Unmarshal(respBytes, &member)
	if err != nil {
		fmt.Printf("with creating struct from bytes : %s", err)
	}
	return member.IDMember
}
