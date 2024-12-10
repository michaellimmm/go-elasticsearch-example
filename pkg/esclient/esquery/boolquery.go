package esquery

import "encoding/json"

type boolQuery struct {
	Must               []QueryType `json:"must,omitempty"`
	MustNot            []QueryType `json:"must_not,omitempty"`
	Should             []QueryType `json:"should,omitempty"`
	Filter             []QueryType `json:"filter,omitempty"`
	MinimumShouldMatch int16       `json:"minimum_should_match,omitempty"`
	Boost              float32     `json:"boost,omitempty"`
}

func (b *boolQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(KeyVal{
		"bool": *b,
	})
}

func (b *boolQuery) SetMust(must ...QueryType) *boolQuery {
	b.Must = append(b.Must, must...)
	return b
}

func (b *boolQuery) SetMustNot(mustNot ...QueryType) *boolQuery {
	b.MustNot = append(b.MustNot, mustNot...)
	return b
}

func (b *boolQuery) SetShould(should ...QueryType) *boolQuery {
	b.Should = append(b.Should, should...)
	return b
}

func (b *boolQuery) SetFilter(filter ...QueryType) *boolQuery {
	b.Filter = append(b.Filter, filter...)
	return b
}

func (b *boolQuery) SetMinimumShouldMatch(min int16) *boolQuery {
	b.MinimumShouldMatch = min
	return b
}

func (b *boolQuery) SetBoost(boost float32) *boolQuery {
	b.Boost = boost
	return b
}

func Bool() *boolQuery {
	return &boolQuery{}
}
