package esclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Index interface {
	CreateIndex(index string, body io.Reader) (*Response[IndexCreationResult], error)
	GetIndeces(index []string, options ...getIndecesOptions) (*Response[map[string]*IndexGetResult], error)
	DeleteIndeces(index []string) (*Response[IndexDeletionResult], error)
}

type IndexCreationResult struct {
	Acknowledged bool   `json:"acknowledged"`
	ShardsAcked  int    `json:"shards_acknowledged"`
	Index        string `json:"index,omitempty"`
}

func (c *client) CreateIndex(index string, body io.Reader) (*Response[IndexCreationResult], error) {
	req, err := http.NewRequest("PUT", c.baseUrl+"/"+index, body)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println(req.Header)

	response := &Response[IndexCreationResult]{
		StatusCode: res.StatusCode,
	}
	response.SetBody(res.Body)

	return response, nil
}

type IndexGetResult struct {
	Aliases  map[string]interface{} `json:"aliases,omitempty"`
	Mapping  map[string]interface{} `json:"mappings,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Warmer   map[string]interface{} `json:"warmer,omitempty"`
}

type getIndecesOptions func(*getIndecesParams)

type getIndecesParams struct {
	features     []string
	httpHeadOnly bool
}

// valid features: ["aliases", "mappings", "settings"]
func GetIndecesWithFeatures(features []string) getIndecesOptions {
	return func(params *getIndecesParams) {
		params.features = features
	}
}

func GetIndecesWithHttpHeadOnly() getIndecesOptions {
	return func(params *getIndecesParams) {
		params.httpHeadOnly = true
	}
}

// Response codes `200`, `404`
// `404` is returned if index does not exist
// `200` is returned if index exists
func (c *client) GetIndeces(index []string, options ...getIndecesOptions) (*Response[map[string]*IndexGetResult], error) {
	params := &getIndecesParams{}
	for _, option := range options {
		option(params)
	}

	uri, err := url.Parse(c.baseUrl + "/" + strings.Join(index, ","))
	if err != nil {
		return nil, err
	}

	if len(params.features) > 0 {
		q := uri.Query()
		q.Add("features", strings.Join(params.features, ","))
		uri.RawQuery = q.Encode()
	}

	method := "GET"
	if params.httpHeadOnly {
		method = "HEAD"
	}

	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &Response[map[string]*IndexGetResult]{
		StatusCode: res.StatusCode,
	}
	response.SetBody(res.Body)

	return response, nil
}

func (c *client) DeleteIndeces(index []string) (*Response[IndexDeletionResult], error) {
	uri, err := url.Parse(c.baseUrl + "/" + strings.Join(index, ","))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", uri.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &Response[IndexDeletionResult]{
		StatusCode: res.StatusCode,
	}
	response.SetBody(res.Body)

	return response, nil
}

type IndexDeletionResult struct {
	Acknowledged bool `json:"acknowledged"`
}
