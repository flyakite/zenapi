package models

import (
	"time"
)

type Link struct {
	Id          int64
	Signal      *Signal    `orm:"rel(fk);null;on_delete(set_null)"`
	Access      []*Access  `orm:"reverse(many)"`
	AccessCount int        `orm:"default(0)"`
	URL         string     `orm:"type(text)"`
	URLHash     string     `orm:"size(20)"`
	UserTrack   *UserTrack `orm:"rel(fk);null;on_delete(set_null)"`

	//last access
	CountryCode string `orm:"size(2);null"` //TODO: country code choices limitation/validation
	City        string `orm:"size(100);null"`
	Device      string `orm:"size(30);null"` //TODO: device choices limitation/validation

	Created  time.Time `orm:"type(datetime);auto_now_add"`
	Modified time.Time `orm:"type(datetime);auto_now"`
}

//multiple fields index
func (l *Link) TableIndex() [][]string {
	return [][]string{
		[]string{"Signal", "URLHash"},
	}
}
