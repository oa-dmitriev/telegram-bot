package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	cmds = map[string]func(*BotRepo, []string) ([]string, error){
		"/def": func(repo *BotRepo, args []string) ([]string, error) {
			return repo.FetchData(strings.Join(args, " "))
		},
	}
)

func isDataLeft(data []string, offset int) bool {
	return len(data) > offset*PAGELEN+PAGELEN
}

func (repo *BotRepo) FetchData(term string) ([]string, error) {
	if repo.RedCon != nil {
		statusCmd := repo.RedCon.Get(term)
		if statusCmd.Err() != nil {
			ans := []string{}
			err := json.Unmarshal([]byte(statusCmd.Val()), &ans)
			if err != nil {
				return nil, err
			}
			log.Printf("\nREDIS FOUND\n%#v\n\n", ans)
			return ans, nil
		}
	}

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

	if repo.RedCon != nil {
		b, err := json.Marshal(ans)
		if err == nil {
			log.Printf("\nREDIS WRITTEN\n%#v\n\n", ans)
			repo.RedCon.Set(term, string(b), time.Hour*24)
		}
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
