package models

import (
	"time"
)

type AccessType int

const (
	LINK_ACCESS AccessType = iota
	OPEN_ACCESS
)

type Access struct {
	Id          int64
	Kind        AccessType `orm:""` //LINK_ACCESS, OPEN_ACCESS
	Signal      *Signal    `orm:"rel(fk);"`
	Receiver    *Receiver  `orm:"rel(fk);null;on_delete(set_null)"` //the receiver who accessed
	UserTrack   *UserTrack `orm:"rel(fk);null;on_delete(set_null)"`
	Link        *Link      `orm:"rel(fk);null;on_delete(set_null)"` //exists if this access is a link access
	IP          string     `orm:"size(40)"`
	CountryCode string     `orm:"size(2);null"` //TODO: country code choices limitation/validation
	City        string     `orm:"size(100);null"`
	Device      string     `orm:"size(30);null"` //TODO: device choices limitation/validation
	UserAgent   *UserAgent `orm:"rel(fk);null;on_delete(set_null)"`
	TzOffset    int        `orm:""`              //in minutes
	Proxy       string     `orm:"size(30);null"` // GoogleImageProxy / Microsoft Proxy,..
	Created     time.Time  `orm:"type(datetime);auto_now_add"`
	Modified    time.Time  `orm:"type(datetime);auto_now"`
}
