package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

// config getting
var config = getConfigList()

// IDMember getting
var IDMember = getMemberID()

func parseInt(s string) (i int) {
	strconv.Atoi(s)
	strconv.ParseInt(s, 10, 64)
	fmt.Sscan(s, &i)

	return i
}

func main() {

	if len(os.Args) > 1 {
		switch param := os.Args[1]; param {

		case "-b":
			boards := getListBoards(IDMember)
			for i := 0; i < len(boards); i++ {
				fmt.Printf("%d %+v\n", i+1, boards[i].Name)
			}
			break

		case "-c":
			if len(os.Args) < 2 {
				fmt.Println("needed number")
				return
			}
			s := os.Args[2]
			cardID := parseInt(s)
			cards := getListCards(IDMember, cardID)
			for i := 0; i < len(cards); i++ {
				fmt.Printf("%d %+v\n", i+1, cards[i].Name)
			}
			break
		}
	}
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

func getListCards(IDMember string, id int) (cards []Card) {
	boards := getListBoards(IDMember)
	fmt.Println("SELECT BOARD: ", string(boards[id-1].Name))

	//"/cards/?limit=10&fields=name&key="
	respBytes := requestTrelloAPI("boards/", string(boards[id-1].ID), "/cards")
	err := json.Unmarshal(respBytes, &cards)

	if err != nil {
		fmt.Printf("with creating struct from bytes : %s", err)
	}

	return cards
}

func getListBoards(IDMember string) (boards []Board) {
	respBytes := requestTrelloAPI("members/", string(IDMember), "/boards")
	err := json.Unmarshal(respBytes, &boards)
	if err != nil {
		fmt.Printf("with creating struct from bytes : %s", err)
	}
	return boards
}

func getMemberID() (id string) {
	respBytes := requestTrelloAPI("tokens/", string(config.ClientKey))

	var member Member
	err := json.Unmarshal(respBytes, &member)
	if err != nil {
		fmt.Printf("with creating struct from bytes : %s", err)
	}
	return member.IDMember
}

func requestTrelloAPI(args ...string) (respBytes []byte) {
	var str strings.Builder
	for i := 0; i < len(args); i++ {
		str.WriteString(args[i])
	}
	var request = strings.Join(
		[]string{
			"https://api.trello.com/1/",
			str.String(),
			"?token=",
			string(config.ClientKey),
			"&key=",
			string(config.AppKey)},
		"")
	respBytes, err := readBody(request)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	return respBytes
}
