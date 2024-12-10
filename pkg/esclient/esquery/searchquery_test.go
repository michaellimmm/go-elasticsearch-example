package esquery_test

import (
	"github/shaolim/go-elasticsearch-example/pkg/esclient/esquery"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchQuery(t *testing.T) {
	tests := []struct {
		expected string
		actual   *esquery.SearchQuery
	}{
		{
			expected: `{
				"size": 10,
				"query": {
					"bool": {
						"should": [
							{
								"match": {
									"additionalProperties.color": {
										"query": "ブルー"
									}
								}
							},
							{
								"term": {
									"additionalProperties.color.keyword": {
										"value": "ブルー"
									}
								}
							}
						]
					}
				}
			}`,
			actual: esquery.NewSearchQueryBuilder().
				SetSize(10).
				SetQuery(
					esquery.Bool().
						SetShould(esquery.Match("additionalProperties.color", "ブルー"),
							esquery.Term("additionalProperties.color.keyword", "ブルー"),
						),
				).
				Build(),
		},
		{
			expected: `{
				"size": 3,
				"query": {
					"bool": {
						"should": [
							{
								"match": {
									"message": {
										"query": "this is a test",
										"boost": 1.0
									}
								}
							}
						]
					}
				},
				"sort": [
					{"record.Updated": {"order": "desc"}},
					{"_score": {"order": "desc"}}
				]
			}`,
			actual: esquery.NewSearchQueryBuilder().
				SetSize(3).
				SetQuery(esquery.Bool().
					SetShould(
						esquery.Match("message", "this is a test").SetBoost(1),
					),
				).
				SetSort(
					esquery.Sort("record.Updated", esquery.OrderDesc),
					esquery.Sort("_score", esquery.OrderDesc),
				).
				Build(),
		},
	}

	for _, test := range tests {
		jsonData, err := test.actual.MarshalJSON()
		assert.Nil(t, err)

		assert.JSONEq(t, test.expected, string(jsonData))
	}
}
