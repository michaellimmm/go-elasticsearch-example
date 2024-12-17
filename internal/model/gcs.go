package model

type GscAttribute struct {
	EventType string `json:"eventType"`
	ObjectID  string `json:"objectId"`
	BucketID  string `json:"bucketId"`
	TenantID  string `json:"tenantID"`
}
