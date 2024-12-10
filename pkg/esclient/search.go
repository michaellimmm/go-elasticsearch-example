package esclient

import (
	"bytes"
	"encoding/json"
	"github/shaolim/go-elasticsearch-example/pkg/esclient/esquery"
	"net/http"
	"reflect"
)

type Search interface {
	Search(index string, query esquery.SearchQuery) (*Response[SearchResult], error)
}

func (c *client) Search(index string, query esquery.SearchQuery) (*Response[SearchResult], error) {
	r, err := query.MarshalJSON()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/"+index+"/_search", bytes.NewReader(r))
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &Response[SearchResult]{
		StatusCode: res.StatusCode,
	}
	response.SetBody(res.Body)

	return response, nil
}

type SearchResult struct {
	TookInMillis    int64         `json:"took,omitempty"`             // search time in milliseconds
	TerminatedEarly bool          `json:"terminated_early,omitempty"` // request terminated early
	ScrollId        string        `json:"_scroll_id,omitempty"`       // only used with Scroll and Scan operations
	Hits            *SearchHits   `json:"hits,omitempty"`             // the actual search hits
	TimedOut        bool          `json:"timed_out,omitempty"`        // true if the search timed out
	Error           *ErrorDetails `json:"error,omitempty"`            // only used in MultiGet
	Shards          *ShardsInfo   `json:"_shards,omitempty"`          // shard information
	Status          int           `json:"status,omitempty"`           // used in MultiSearch
	PitId           string        `json:"pit_id,omitempty"`           // Point In Time ID
}

func (r *SearchResult) TotalHits() int64 {
	if r != nil && r.Hits != nil && r.Hits.TotalHits != nil {
		return r.Hits.TotalHits.Value
	}
	return 0
}

func (r *SearchResult) Each(typ reflect.Type) []interface{} {
	if r.Hits == nil || r.Hits.Hits == nil || len(r.Hits.Hits) == 0 {
		return nil
	}
	slice := make([]interface{}, 0, len(r.Hits.Hits))
	for _, hit := range r.Hits.Hits {
		v := reflect.New(typ).Elem()
		if hit.Source == nil {
			slice = append(slice, v.Interface())
			continue
		}
		if err := json.Unmarshal(hit.Source, v.Addr().Interface()); err == nil {
			slice = append(slice, v.Interface())
		}
	}
	return slice
}

type SearchHits struct {
	TotalHits *TotalHits   `json:"total,omitempty"`
	MaxScore  *float64     `json:"max_score,omitempty"`
	Hits      []*SearchHit `json:"hits,omitempty"`
}

type TotalHits struct {
	Value    int64  `json:"value"`    // value of the total hit count
	Relation string `json:"relation"` // how the value should be interpreted: accurate ("eq") or a lower bound ("gte")
}

type SearchHit struct {
	Score  *float64        `json:"_score,omitempty"`  // computed score
	Index  string          `json:"_index,omitempty"`  // index name
	Id     string          `json:"_id,omitempty"`     // external or internal
	Sort   []interface{}   `json:"sort,omitempty"`    // sort information
	Source json.RawMessage `json:"_source,omitempty"` // stored document source
}
