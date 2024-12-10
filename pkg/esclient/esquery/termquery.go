package esquery

import "encoding/json"

type termQuery struct {
	Field           string   `json:"-"`
	Value           string   `json:"value"`
	Boost           *float32 `json:"boost,omitempty"`
	CaseInsensitive *bool    `json:"case_insensitive,omitempty"`
}

func (t *termQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(KeyVal{
		"term": KeyVal{
			t.Field: *t,
		},
	})
}

func (t *termQuery) SetBoost(boost float32) *termQuery {
	t.Boost = &boost
	return t
}

func (t *termQuery) SetCaseInsensitive(isCaseSensitive bool) *termQuery {
	t.CaseInsensitive = &isCaseSensitive
	return t
}

func Term(field string, value string) *termQuery {
	return &termQuery{
		Field: field,
		Value: value,
	}
}
