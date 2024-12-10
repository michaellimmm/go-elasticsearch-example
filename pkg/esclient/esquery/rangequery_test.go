package esquery_test

import (
	"encoding/json"
	"github/shaolim/go-elasticsearch-example/pkg/esclient/esquery"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeQuery(t *testing.T) {
	expected := `{
		"range": {
			"fake_field": {
				"gt": 10,
				"gte": 20,
				"lt": 30,
				"lte": 40,
				"format": "strict_date_optional_time",
				"relation": "INTERSECTS",
				"time_zone": "UTC",
				"boost": 1.5
			}
		}
	}`

	actual := esquery.Range("fake_field").
		SetGt(10).
		SetGte(20).
		SetLt(30).
		SetLte(40).
		SetFormat("strict_date_optional_time").
		SetRelation(esquery.INTERSECTS).
		SetTimezone("UTC").
		SetBoost(1.5)

	jsonData, err := json.Marshal(actual)
	assert.Nil(t, err)

	assert.JSONEq(t, expected, string(jsonData))
}
