package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	apiURL = url.URL{
		Scheme: "https",
		Host:   "api.urbandictionary.com",
		Path:   "v0/define",
	}
	PAGELEN = 3
)

type DefinitionData struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type Data struct {
	Definitions []*DefinitionData `json:"list"`
}

var (
	cmds = map[string]func([]string) ([]string, error){
		"/def": func(args []string) ([]string, error) {
			return FetchData(strings.Join(args, " "))
		},
	}
)

func isDataLeft(data []string, offset int) bool {
	return len(data) > offset*PAGELEN+PAGELEN
}

func FetchData(term string) ([]string, error) {
	u := apiURL
	p := url.Values{
		"term": []string{term},
	}
	u.RawQuery = p.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := Data{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	ans := make([]string, len(data.Definitions))
	for i := range data.Definitions {
		ans[i] = data.Definitions[i].Definition
	}
	return ans, nil
}

func GetPage(data []string, offset int) []string {
	start := offset * PAGELEN
	if start > len(data) {
		return nil
	}
	end := offset*PAGELEN + PAGELEN
	if len(data) < end {
		end = len(data)
	}
	return data[start:end]
}
