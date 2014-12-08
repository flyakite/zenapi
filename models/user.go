package models

import ()

type User struct {
	Id       int    `orm:"column(id);auto"`
	Username string `orm:"column(username);size(100)"`
	Email    string `orm:column(email);size(255)`
}
