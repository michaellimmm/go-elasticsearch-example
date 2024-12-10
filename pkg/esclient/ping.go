package esclient

import (
	"net/http"
)

type Ping interface {
	Ping(options ...pingParamsOptions) (*Response[PingResult], error)
}

type PingResult struct {
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Name        string `json:"name"`
	TagLine     string `json:"tagline"`
	Version     struct {
		BuildDate           string `json:"build_date"`
		BuildFlavor         string `json:"build_flavor"`
		BuildHash           string `json:"build_hash"`
		BuildSnapshot       bool   `json:"build_snapshot"`
		BuildType           string `json:"build_type"`
		LuceneVersion       string `json:"lucene_version"`
		MinimumIndexVersion string `json:"minimum_index_compatibility_version"`
		MinimumWireVersion  string `json:"minimum_wire_compatibility_version"`
		Number              string `json:"number"`
	} `json:"version"`
}

type pingParamsOptions func(*pingParams)

type pingParams struct {
	httpHeadOnly bool
}

func PingWithHttpHeadOnly() func(*pingParams) {
	return func(p *pingParams) {
		p.httpHeadOnly = true
	}
}

func (c *client) Ping(options ...pingParamsOptions) (*Response[PingResult], error) {
	params := &pingParams{}
	for _, option := range options {
		option(params)
	}

	method := "GET"
	if params.httpHeadOnly {
		method = "HEAD"
	}

	req, err := http.NewRequest(method, c.baseUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &Response[PingResult]{
		StatusCode: res.StatusCode,
	}
	response.SetBody(res.Body)

	return response, nil
}
