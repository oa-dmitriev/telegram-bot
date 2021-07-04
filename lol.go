package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var dataUrl = "https://wordsapiv1.p.mashape.com/words"

type DefinitionData struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type Data struct {
	Definitions []*DefinitionData `json:"list"`
}

func main() {
	//	url := "https://www.wordsapi.com/mashape/words/star/definitions?encrypted=8cfdb18be722949bea9207bded58bdbbaeb5290931fd95b8"
	//	url := fmt.Sprintf("%s/%s/definitions?encrypted=encrypted=8cfdb18be722949bea9207bded58bdbbaeb5290931fd95b8", dataUrl, "word")
	u := url.URL{
		Scheme: "https",
		Host:   "api.urbandictionary.com",
		Path:   "v0/define",
	}
	args := []string{"black", "out"}
	p := url.Values{
		"term": []string{strings.Join(args, " ")},
	}
	u.RawQuery = p.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal("http get error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read error: ", err)
	}
	wd := Data{}
	err = json.Unmarshal(body, &wd)
	if err != nil {
		log.Fatal("unmarhsalling error: ", err)
	}
	log.Printf("word: %s\ndef: %s\n", wd.Definitions[0].Word, wd.Definitions[0].Definition)
}
