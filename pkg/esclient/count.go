package esclient

import (
	"bytes"
	"encoding/json"
	"github/shaolim/kakashi/pkg/esclient/esquery"
	"net/http"
)

type Count interface {
	Count(index string, query esquery.QueryType) (*Response[CountResponse], error)
}

func (c *client) Count(index string, query esquery.QueryType) (*Response[CountResponse], error) {
	queryRequest := esquery.KeyVal{
		"query": query,
	}

	r, err := json.Marshal(queryRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/"+index+"/_count", bytes.NewReader(r))
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &Response[CountResponse]{
		StatusCode: res.StatusCode,
	}
	response.SetBody(res.Body)

	return response, nil
}

type CountResponse struct {
	Count           int64       `json:"count"`
	TerminatedEarly bool        `json:"terminated_early,omitempty"`
	Shards          *ShardsInfo `json:"_shards,omitempty"`
}
