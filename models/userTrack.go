package models

import ()

type UserTrack struct {
	Id int64
	// Name  string `orm:"size(200)"`
	// Email string `orm:"size(255)"`
	Ass string `orm:"size(32)"`
}
