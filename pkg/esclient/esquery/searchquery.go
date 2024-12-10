package esquery

import "encoding/json"

type KeyVal map[string]interface{}

type QueryType interface {
	json.Marshaler
}

type SearchQuery struct {
	Size  uint32    `json:"size,omitempty"`
	Query QueryType `json:"query,omitempty"`
	From  uint32    `json:"from,omitempty"`
	Sort  []*sort   `json:"sort,omitempty"`
}

func (s *SearchQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

type sort struct {
	Field string
	Order Order
}

func Sort(field string, order Order) *sort {
	return &sort{field, order}
}

func (s *sort) MarshalJSON() ([]byte, error) {
	return json.Marshal(KeyVal{
		s.Field: KeyVal{
			"order": s.Order,
		},
	})
}

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type SearchQueryBuilder struct {
	searchQuery *SearchQuery
}

func NewSearchQueryBuilder() *SearchQueryBuilder {
	return &SearchQueryBuilder{searchQuery: &SearchQuery{}}
}

func (s *SearchQueryBuilder) SetSize(size uint32) *SearchQueryBuilder {
	s.searchQuery.Size = size
	return s
}

func (s *SearchQueryBuilder) SetFrom(from uint32) *SearchQueryBuilder {
	s.searchQuery.From = from
	return s
}

func (s *SearchQueryBuilder) SetQuery(query QueryType) *SearchQueryBuilder {
	s.searchQuery.Query = query
	return s
}

func (s *SearchQueryBuilder) SetSort(sort ...*sort) *SearchQueryBuilder {
	s.searchQuery.Sort = append(s.searchQuery.Sort, sort...)
	return s
}

func (s *SearchQueryBuilder) Build() *SearchQuery {
	return s.searchQuery
}
