package domain

import (
	"fmt"
	"strings"

	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

func ConvertDBVocabToDomainData(dbVocab []*repository.DBVocabulary) []*DefinitionData {
	res := make([]*DefinitionData, 0, len(dbVocab))
	for _, v := range dbVocab {
		res = append(res, &DefinitionData{
			Word:       v.Word,
			Definition: v.Definition,
		})
	}
	return res
}

func ToString(domainData []*DefinitionData) string {
	if len(domainData) == 0 {
		return ""
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "*%s* - %s", domainData[0].Word, domainData[0].Definition)
	for i := 1; i < len(domainData); i++ {
		fmt.Fprintf(&sb, "\n\n*%s* - %s", domainData[i].Word, domainData[i].Definition)
	}
	return sb.String()
}

func ConvertJokes(jokes []*repository.DBJoke) []*JokeInfo {
	res := make([]*JokeInfo, 0, len(jokes))
	for _, joke := range jokes {
		res = append(res, &JokeInfo{
			Category: joke.Category,
			Type:     joke.Type,
			Joke:     joke.Joke,
			Setup:    joke.Setup,
			Delivery: joke.Delivery,
			ID:       joke.ID,
		})
	}
	return res
}

func JokesToString(jokes []*JokeInfo) string {
	var sb strings.Builder

	for _, joke := range jokes {
		fmt.Fprintln(&sb, jokeToSring(joke))
	}
	return sb.String()
}

func jokeToSring(joke *JokeInfo) string {
	if joke.Type == "single" {
		return joke.Joke
	}
	if joke.Type == "twopart" {
		return fmt.Sprintf("%s\n%s", joke.Setup, joke.Delivery)
	}
	return "Колобок повесился"
}
