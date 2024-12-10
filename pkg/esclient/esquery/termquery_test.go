package esquery_test

import (
	"encoding/json"
	"github/shaolim/go-elasticsearch-example/pkg/esclient/esquery"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTermQuery(t *testing.T) {
	expected := `{
		"term": {
            "user.id": {
                "value": "kimchy",
                "boost": 1,
				"case_insensitive": false
            }
        }
	}`

	actual := esquery.Term("user.id", "kimchy").
		SetBoost(1).
		SetCaseInsensitive(false)

	jsonData, err := json.Marshal(actual)
	assert.Nil(t, err)

	assert.JSONEq(t, expected, string(jsonData))
}
