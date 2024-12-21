package esclient

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BulkableRequest interface {
	String() (string, error)
}

type BulkRequests struct {
	BulkRequests []BulkableRequest
}

func (b *BulkRequests) Add(bulkRequest BulkableRequest) *BulkRequests {
	b.BulkRequests = append(b.BulkRequests, bulkRequest)
	return b
}

func (b *BulkRequests) Length() int {
	return len(b.BulkRequests)
}

func (b BulkRequests) String() (string, error) {
	var sb strings.Builder
	for _, bulkRequest := range b.BulkRequests {
		str, err := bulkRequest.String()
		if err != nil {
			return "", err
		}
		sb.WriteString(str)
	}
	return sb.String(), nil
}

type bulkDeleteRequest struct {
	Index string `json:"_index,omitempty"`
	Id    string `json:"_id"`
}

func NewBulkDeleteRequest(id string) *bulkDeleteRequest {
	return &bulkDeleteRequest{
		Id: id,
	}
}

func (b *bulkDeleteRequest) SetIndex(index string) *bulkDeleteRequest {
	b.Index = index
	return b
}

func (b *bulkDeleteRequest) String() (string, error) {
	p, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{\"delete\":%v}\n", string(p)), nil
}

type bulkIndexRequest struct {
	Index         string      `json:"_index,omitempty"`
	Id            string      `json:"_id"`
	Routing       string      `json:"routing,omitempty"`
	Pipeline      string      `json:"pipeline,omitempty"`
	IfSeqNo       int64       `json:"if_seq_no,omitempty"`
	IfPrimaryTerm int64       `json:"if_primary_term,omitempty"`
	Doc           interface{} `json:"-"`
}

func NewBulkIndexRequest() *bulkIndexRequest {
	return &bulkIndexRequest{}
}

func (b *bulkIndexRequest) SetIndex(index string) *bulkIndexRequest {
	b.Index = index
	return b
}

func (b *bulkIndexRequest) SetId(id string) *bulkIndexRequest {
	b.Id = id
	return b
}

func (b *bulkIndexRequest) SetRouting(routing string) *bulkIndexRequest {
	b.Routing = routing
	return b
}

func (b *bulkIndexRequest) SetPipeline(pipeline string) *bulkIndexRequest {
	b.Pipeline = pipeline
	return b
}

func (b *bulkIndexRequest) SetIfSeqNo(ifSeqNo int64) *bulkIndexRequest {
	b.IfSeqNo = ifSeqNo
	return b
}

func (b *bulkIndexRequest) SetIfPrimaryTerm(ifPrimaryTerm int64) *bulkIndexRequest {
	b.IfPrimaryTerm = ifPrimaryTerm
	return b
}

func (b *bulkIndexRequest) SetDoc(doc interface{}) *bulkIndexRequest {
	b.Doc = doc
	return b
}

func (b *bulkIndexRequest) String() (string, error) {
	p, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	action := fmt.Sprintf("{\"index\":%v}", string(p))

	doc := "{}"
	if b.Doc != nil {
		_docs, err := json.Marshal(b.Doc)
		if err != nil {
			return "", err
		}
		doc = string(_docs)
	}

	return fmt.Sprintf("%s\n%s\n", action, doc), nil
}

type bulkUpdateRequest struct {
	Index           string      `json:"_index,omitempty"`
	Id              string      `json:"_id"`
	RetryOnConflict int         `json:"retry_on_conflict,omitempty"`
	Routing         string      `json:"routing,omitempty"`
	IfSeqNo         int64       `json:"if_seq_no,omitempty"`
	IfPrimaryTerm   int64       `json:"if_primary_term,omitempty"`
	Doc             interface{} `json:"-"`
}

func NewBulkUpdateRequest(id string) *bulkUpdateRequest {
	return &bulkUpdateRequest{
		Id: id,
	}
}

func (b *bulkUpdateRequest) SetIndex(index string) *bulkUpdateRequest {
	b.Index = index
	return b
}

func (b *bulkUpdateRequest) SetRetryOnConflict(retryOnConflict int) *bulkUpdateRequest {
	b.RetryOnConflict = retryOnConflict
	return b
}

func (b *bulkUpdateRequest) SetRouting(routing string) *bulkUpdateRequest {
	b.Routing = routing
	return b
}

func (b *bulkUpdateRequest) SetIfSeqNo(ifSeqNo int64) *bulkUpdateRequest {
	b.IfSeqNo = ifSeqNo
	return b
}

func (b *bulkUpdateRequest) SetIfPrimaryTerm(ifPrimaryTerm int64) *bulkUpdateRequest {
	b.IfPrimaryTerm = ifPrimaryTerm
	return b
}

func (b *bulkUpdateRequest) SetDoc(doc interface{}) *bulkUpdateRequest {
	b.Doc = doc
	return b
}

func (b *bulkUpdateRequest) String() (string, error) {
	p, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	action := fmt.Sprintf("{\"update\":%v}", string(p))

	doc := "{}"
	if b.Doc != nil {
		_docs, err := json.Marshal(b.Doc)
		if err != nil {
			return "", err
		}
		doc = string(_docs)
	}

	return fmt.Sprintf("%s\n{\"doc\":%s}\n", action, doc), nil
}

type bulkCreateRequest struct {
	Index    string      `json:"_index,omitempty"`
	Id       string      `json:"_id"`
	Routing  string      `json:"routing,omitempty"`
	Pipeline string      `json:"pipeline,omitempty"`
	Doc      interface{} `json:"-"`
}

func NewBulkCreateRequest(id string) *bulkCreateRequest {
	return &bulkCreateRequest{
		Id: id,
	}
}

func (b *bulkCreateRequest) SetIndex(index string) *bulkCreateRequest {
	b.Index = index
	return b
}

func (b *bulkCreateRequest) SetRouting(routing string) *bulkCreateRequest {
	b.Routing = routing
	return b
}

func (b *bulkCreateRequest) SetPipeline(pipeline string) *bulkCreateRequest {
	b.Pipeline = pipeline
	return b
}

func (b *bulkCreateRequest) SetDoc(doc interface{}) *bulkCreateRequest {
	b.Doc = doc
	return b
}

func (b *bulkCreateRequest) String() (string, error) {
	p, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	action := fmt.Sprintf("{\"create\":%v}", string(p))

	doc := "{}"
	if b.Doc != nil {
		_docs, err := json.Marshal(b.Doc)
		if err != nil {
			return "", err
		}
		doc = string(_docs)
	}

	return fmt.Sprintf("%s\n%s\n", action, doc), nil
}
