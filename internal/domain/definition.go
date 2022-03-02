package domain

type DefinitionData struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
	ThumbsUp   int32  `json:"thumbs_up"`
	ThumbDown  int32  `json:"thumbs_down"`
}

type Data struct {
	Definitions []*DefinitionData `json:"list"`
}
