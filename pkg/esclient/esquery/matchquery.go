package esquery

import "encoding/json"

type matchQuery struct {
	Field              string   `json:"-"`
	Query              string   `json:"query"`
	Boost              *float64 `json:"boost,omitempty"`
	MinimumShouldMatch *int     `json:"minimum_should_match,omitempty"`
	Fuzziness          *float64 `json:"fuzziness,omitempty"`
	PrefixLength       *int     `json:"prefix_length,omitempty"`
	MaxExpansions      *int     `json:"max_expansions,omitempty"`
}

func (m *matchQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(KeyVal{
		"match": KeyVal{
			m.Field: *m,
		},
	})
}

func (m *matchQuery) SetBoost(boost float64) *matchQuery {
	m.Boost = &boost
	return m
}

func (m *matchQuery) SetMinimumShouldMatch(min int) *matchQuery {
	m.MinimumShouldMatch = &min
	return m
}

func (m *matchQuery) SetFuzziness(fuzziness float64) *matchQuery {
	m.Fuzziness = &fuzziness
	return m
}

func (m *matchQuery) SetPrefixLength(prefixLength int) *matchQuery {
	m.PrefixLength = &prefixLength
	return m
}

func (m *matchQuery) SetMaxExpansions(maxExpansions int) *matchQuery {
	m.MaxExpansions = &maxExpansions
	return m
}

func Match(field, query string) *matchQuery {
	return &matchQuery{
		Field: field,
		Query: query,
	}
}

type matchAllQuery struct {
	Boost *float64 `json:"boost,omitempty"`
}

func MatchAll() *matchAllQuery {
	return &matchAllQuery{}
}

func (m *matchAllQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(KeyVal{
		"match_all": *m,
	})
}

func (m *matchAllQuery) SetBoost(boost float64) *matchAllQuery {
	m.Boost = &boost
	return m
}
