package esclient_test

import (
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBulkRequest(t *testing.T) {
	expected := `{"update":{"_index":"products","_id":"1"}}
{"doc":{"price":649.99}}
{"create":{"_index":"products","_id":"2"}}
{"name":"Laptop","price":1299.99,"in_stock":false}
{"delete":{"_index":"products","_id":"1"}}
{"index":{"_index":"products","_id":"3"}}
{"name":"Laptop","price":1299.99,"in_stock":false}
`

	type product struct {
		Name    string  `json:"name,omitempty"`
		Price   float32 `json:"price,omitempty"`
		InStock *bool   `json:"in_stock,omitempty"`
	}

	actual, err := esclient.BulkRequests{
		BulkRequests: []esclient.BulkableRequest{
			esclient.NewBulkUpdateRequest("1").SetIndex("products").
				SetDoc(product{Price: 649.99}),
			esclient.NewBulkCreateRequest("2").SetIndex("products").
				SetDoc(product{Name: "Laptop", InStock: boolPtr(false), Price: 1299.99}),
			esclient.NewBulkDeleteRequest("1").SetIndex("products"),
			esclient.NewBulkIndexRequest().SetId("3").SetIndex("products").
				SetDoc(product{Name: "Laptop", InStock: boolPtr(false), Price: 1299.99})},
	}.
		String()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func boolPtr(b bool) *bool {
	return &b
}
