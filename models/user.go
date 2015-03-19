package models

import ()

type User struct {
	Id       int    `orm:"auto"`
	Username string `orm:"size(100)"`
	Email    string `orm:"size(255)"`
}
