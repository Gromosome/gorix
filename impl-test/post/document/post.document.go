package document

import "time"

type Post struct {
	ID string `json:"id" bson:"_id,omitempty" couchdb:"_id,omitempty" document:"id"`

	// CouchDB uses this. MongoDB ignores it.
	Rev string `json:"rev,omitempty" bson:"-" couchdb:"_rev,omitempty" document:"rev"`

	Title     string    `json:"title" bson:"title" couchdb:"title" document:"field,index"`
	Slug      string    `json:"slug" bson:"slug" couchdb:"slug" document:"field,uniqueIndex"`
	Content   string    `json:"content" bson:"content" couchdb:"content" document:"field"`
	Tags      []string  `json:"tags" bson:"tags" couchdb:"tags" document:"field,index"`
	Published bool      `json:"published" bson:"published" couchdb:"published" document:"field,index"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt" couchdb:"createdAt" document:"field,index"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt" couchdb:"updatedAt" document:"field"`
}

func (Post) CollectionName() string {
	return "posts"
}
