package esquery_test

import (
	"encoding/json"
	"github/shaolim/kakashi/pkg/esclient/esquery"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolQuery(t *testing.T) {
	expected := `{
		"bool" : {
			"must" : [
				{
					"term" : { 
						"user.id" : { 
							"value" : "kimchy"
						} 
					}
				}
			],
			"filter": [
				{
					"term" : { 
						"tags" : { 
							"value" : "production"
						} 
					}
				}
			],
			"must_not" : [
				{
					"range" : {
						"age" : { "gte" : 10, "lte" : 20 }
					}
				}
			],
			"should" : [
				{ 
					"term" : { 
						"tags" : { 
							"value" : "env1"
						} 
					} 
				},
				{ 
					"term" : { 
						"tags" : { 
							"value" : "deployed"
						} 
					} 
				}
			],
			"minimum_should_match" : 1,
			"boost" : 1.0
		}
	}`

	actual := esquery.Bool().
		SetMust(
			esquery.Term("user.id", "kimchy"),
		).
		SetFilter(
			esquery.Term("tags", "production"),
		).
		SetMustNot(
			esquery.Range("age").SetGte(10).SetLte(20),
		).
		SetShould(
			esquery.Term("tags", "env1"),
			esquery.Term("tags", "deployed"),
		).
		SetMinimumShouldMatch(1).
		SetBoost(1.0)

	jsonData, err := json.Marshal(actual)
	assert.Nil(t, err)

	assert.JSONEq(t, expected, string(jsonData))
}
