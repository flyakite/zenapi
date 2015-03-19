package models

import ()

type Receiver struct {
	Id    int64
	Email string `orm:"size(255)";pk`
	Names string `orm:"type(text)`
}
