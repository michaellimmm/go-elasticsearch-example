package esquery_test

import (
	"encoding/json"
	"github/shaolim/kakashi/pkg/esclient/esquery"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchQuery(t *testing.T) {
	expected := `{
		"match": {
            "message": {
                "query": "this is a test",
				"boost": 1
            }
        }
	}`

	actual := esquery.Match("message", "this is a test").
		SetBoost(1)

	jsonData, err := json.Marshal(actual)
	assert.Nil(t, err)

	assert.JSONEq(t, expected, string(jsonData))
}

func TestMatchAllQuery(t *testing.T) {
	expected := `{
		"match_all": {}
	}`

	actual := esquery.MatchAll()

	jsonData, err := json.Marshal(actual)
	assert.Nil(t, err)

	assert.JSONEq(t, expected, string(jsonData))
}
