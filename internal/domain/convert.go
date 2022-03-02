package domain

import "github.com/oa-dmitriev/telegram-bot/internal/repository"

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
