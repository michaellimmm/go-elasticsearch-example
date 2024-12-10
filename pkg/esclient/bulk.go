package esclient

import (
	"net/http"
	"strings"
)

type Bulk interface {
	Bulk(index string, bulkRequest BulkableRequest) (*Response[BulkResult], error)
}

func (c *client) Bulk(index string, bulkRequest BulkableRequest) (*Response[BulkResult], error) {
	r, err := bulkRequest.String()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseUrl+"/"+index+"/_bulk", strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-ndjson")

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &Response[BulkResult]{
		StatusCode: res.StatusCode,
	}
	response.SetBody(res.Body)

	return response, nil
}

type BulkResult struct {
	Took   int                            `json:"took,omitempty"`
	Errors bool                           `json:"errors,omitempty"`
	Items  []map[string]*BulkResponseItem `json:"items,omitempty"`
}

type BulkResponseItem struct {
	Index       string        `json:"_index,omitempty"`
	Id          string        `json:"_id,omitempty"`
	Version     int64         `json:"_version,omitempty"`
	Result      string        `json:"result,omitempty"`
	Shards      *ShardsInfo   `json:"_shards,omitempty"`
	PrimaryTerm int64         `json:"_primary_term,omitempty"`
	Status      int           `json:"status,omitempty"`
	Error       *ErrorDetails `json:"error,omitempty"`
}

func (r *BulkResult) Indexed() []*BulkResponseItem {
	return r.ByAction("index")
}

func (r *BulkResult) Created() []*BulkResponseItem {
	return r.ByAction("create")
}

func (r *BulkResult) Updated() []*BulkResponseItem {
	return r.ByAction("update")
}

func (r *BulkResult) Deleted() []*BulkResponseItem {
	return r.ByAction("delete")
}

func (r *BulkResult) ByAction(action string) []*BulkResponseItem {
	if r.Items == nil {
		return nil
	}
	var items []*BulkResponseItem
	for _, item := range r.Items {
		if result, found := item[action]; found {
			items = append(items, result)
		}
	}
	return items
}

func (r *BulkResult) ById(id string) []*BulkResponseItem {
	if r.Items == nil {
		return nil
	}
	var items []*BulkResponseItem
	for _, item := range r.Items {
		for _, result := range item {
			if result.Id == id {
				items = append(items, result)
			}
		}
	}
	return items
}

func (r *BulkResult) Failed() []*BulkResponseItem {
	if r.Items == nil {
		return nil
	}
	var errors []*BulkResponseItem
	for _, item := range r.Items {
		for _, result := range item {
			if !(result.Status >= 200 && result.Status <= 299) {
				errors = append(errors, result)
			}
		}
	}
	return errors
}

func (r *BulkResult) Succeeded() []*BulkResponseItem {
	if r.Items == nil {
		return nil
	}
	var succeeded []*BulkResponseItem
	for _, item := range r.Items {
		for _, result := range item {
			if result.Status >= 200 && result.Status <= 299 {
				succeeded = append(succeeded, result)
			}
		}
	}
	return succeeded
}

type ShardsInfo struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
	Skipped    int `json:"skipped"`
}

type ErrorDetails struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	IndexUUID string `json:"index_uuid"`
	Index     string `json:"index"`
	Shard     int    `json:"shard"`
}
