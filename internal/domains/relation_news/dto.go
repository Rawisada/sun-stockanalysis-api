package relation_news

type CreateRelationNewsInput struct {
	Body struct {
		Symbol          string   `json:"symbol"`
		RelationSymbols []string `json:"relation_symbols"`
	}
}
