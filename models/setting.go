package models

import ()

type Setting struct {
	Id                int    `orm:"pk;auto"`
	Email             string `orm:"size(255)"`
	TrackByDefault    bool   `orm:"null"`
	IsNotifyByEmail   bool   `orm:"null"`
	IsNotifyByDesktop bool   `orm:"null"`
	IsDailyReport     bool   `orm:"null"`
	IsWeeklyReport    bool   `orm:"null"`
}
