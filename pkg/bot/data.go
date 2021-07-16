package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func isDataLeft(data []string, offset int) bool {
	return len(data) > offset*PAGELEN+PAGELEN
}

func (repo *BotRepo) FetchData(term string) ([]string, error) {
	if repo.redClient != nil {
		statusCmd := repo.redClient.Get(term)
		if statusCmd.Err() == nil {
			ans := []string{}
			err := json.Unmarshal([]byte(statusCmd.Val()), &ans)
			if err != nil {
				log.Println("\nERROR: ", err)
				return nil, err
			}
			log.Printf("\nFOUND IN REDIS%#v\n", ans)
			return ans, nil
		}
		log.Println("\nNO REDIS ENTRY")
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

	if repo.redClient != nil {
		b, err := json.Marshal(ans)
		if err == nil {
			statusCmd := repo.redClient.Set(term, string(b), time.Hour*24)
			if statusCmd.Err() != nil {
				log.Printf("\nREDIS WRITTEN ERROR\n%s\n\n", statusCmd.Err().Error())
			} else {
				log.Printf("\nREDIS WRITTEN\n%#v\n\n", ans)
			}
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
