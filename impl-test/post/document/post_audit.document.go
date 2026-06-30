package document

import "time"

type PostAudit struct {
	ID string `json:"id" bson:"_id,omitempty" couchdb:"_id,omitempty" document:"id"`

	Rev string `json:"rev,omitempty" bson:"-" couchdb:"_rev,omitempty" document:"rev"`

	PostID    string    `json:"postId" bson:"postId" couchdb:"postId" document:"field,index"`
	Action    string    `json:"action" bson:"action" couchdb:"action" document:"field,index"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt" couchdb:"createdAt" document:"field,index"`
}

func (PostAudit) CollectionName() string {
	return "post_audits"
}
