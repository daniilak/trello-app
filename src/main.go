package main

import (
	"encoding/json"
	"flag"
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
	ID               string
	Closed           bool
	DateLastActivity string
	Desc             string
	IDBoard          string
	IDList           string
	Name             string
	Due              string
}

// Member structure
type Member struct {
	ID          string
	IDMember    string
	Identifier  string
	dateCreated string
	dateExpires string
}

// Board structure
type Board struct {
	ID               string
	Desc             string
	DescData         string
	Closed           bool
	IDOrganization   string
	ShortLink        string
	DateLastActivity string
	URL              string
	Name             string
	Starred          bool
}

// StringFlag structure
type stringFlag struct {
	set      bool
	value    string
	typeflag string
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

var cFlag, bFlag stringFlag

func parseInt(s string) (i int) {
	strconv.Atoi(s)
	strconv.ParseInt(s, 10, 64)
	fmt.Sscan(s, &i)

	return i
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func parseArgs() bool {
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println("	-b 	         get list boards")
		fmt.Println("	-c <number>	 get list cards")
		return true
	}
	return false
}

func parseFlags() {
	flag.Var(&bFlag, "b", "the list of boards")
	flag.Var(&cFlag, "c", "the cards of board")

	flag.Parse()

	if cFlag.set {
		printCards(getListCards(IDMember, cFlag.value))
	}
	if bFlag.set {
		printBoards(getListBoards(IDMember))
	}
}

func main() {
	if len(os.Args) > 1 {
		if !parseArgs() {
			parseFlags()
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

func getListCards(IDMember string, id string) (cards []Card) {
	boards := getListBoards(IDMember)

	fmt.Println("SELECT BOARD: ", string(boards[parseInt(id)-1].Name))

	respBytes := requestTrelloAPI("boards/", string(boards[parseInt(id)-1].ID), "/cards")
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
	//fmt.Println(request)
	respBytes, err := readBody(request)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	return respBytes
}

func printBoards(boards []Board) {
	for i := 0; i < len(boards); i++ {
		fmt.Printf("%d %+v\n", i+1, boards[i].Name)
	}
}

func printCards(cards []Card) {
	for i := 0; i < len(cards); i++ {
		fmt.Printf("%d %+v %+v\n", i+1, cards[i].Name, cards[i].Due)
	}
}
