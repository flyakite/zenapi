package models

type UserAgent struct {
	Id   int    `orm:"auto"`
	Text string `orm:"type(text)"`
	Hash string `orm:"size(32)"`
}
