package esquery

import "encoding/json"

type rangeQuery struct {
	Field    string      `json:"-"`
	Gt       interface{} `json:"gt,omitempty"`
	Gte      interface{} `json:"gte,omitempty"`
	Lt       interface{} `json:"lt,omitempty"`
	Lte      interface{} `json:"lte,omitempty"`
	Boost    *float32    `json:"boost,omitempty"`
	Format   string      `json:"format,omitempty"`
	TimeZone string      `json:"time_zone,omitempty"`
	Relation Relation    `json:"relation,omitempty"`
}

type Relation string

const (
	INTERSECTS Relation = "INTERSECTS"
	CONTAINS   Relation = "CONTAINS"
	WITHIN     Relation = "WITHIN"
)

func (r *rangeQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(KeyVal{
		"range": KeyVal{
			r.Field: *r,
		},
	})
}

func (r *rangeQuery) SetGt(gt interface{}) *rangeQuery {
	r.Gt = gt
	return r
}

func (r *rangeQuery) SetGte(gte interface{}) *rangeQuery {
	r.Gte = gte
	return r
}

func (r *rangeQuery) SetLt(lt interface{}) *rangeQuery {
	r.Lt = lt
	return r
}

func (r *rangeQuery) SetLte(lte interface{}) *rangeQuery {
	r.Lte = lte
	return r
}

func (r *rangeQuery) SetBoost(boost float32) *rangeQuery {
	r.Boost = &boost
	return r
}

func (r *rangeQuery) SetFormat(format string) *rangeQuery {
	r.Format = format
	return r
}

func (r *rangeQuery) SetTimezone(timezone string) *rangeQuery {
	r.TimeZone = timezone
	return r
}

func (r *rangeQuery) SetRelation(relation Relation) *rangeQuery {
	r.Relation = relation
	return r
}

func Range(field string) *rangeQuery {
	return &rangeQuery{Field: field}
}
