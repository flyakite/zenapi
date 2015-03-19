/*
Signal
ref: http://beego.me/docs/mvc/model/models.md
*/
package models

import (
	"time"
)

type NotificationSettingType int

const (
	EMAIL_ONLY NotificationSettingType = iota
	DESKTOP_ONLY
	NOTIFICATION_DISABLED
)

type EmailClientType int

const (
	GMAIL_CLIENT EmailClientType = iota
	OUTLOOK_CLIENT
)

var SIGNAL_RECEIVER_KEYS = [3]string{"to", "cc", "bcc"}

const SIGNAL_RECEIVER_EMAILS_SEPERATOR = ";"

type Signal struct {
	Id              int64
	Token           string    `orm:"size(50)";index;unique`
	Email           string    `orm:"size(255);index"`
	Subject         string    `orm:"type(text)"`
	Body            string    `orm:"type(text)"`
	Accesses        []*Access `orm:"reverse(many)"`
	Links           []*Link   `orm:"reverse(many)"`
	AccessCount     int       `orm:"default(0)"`
	ReceiverArchive string    `orm:"type(text);null"`   //seperated by ";"
	ReceiverEmails  string    `orm:"type(text)"`        //seperated by SIGNAL_RECEIVER_EMAILS_SEPERATOR ";"
	TzOffset        int       `orm:"column(tz_offset)"` //in minutes

	/*
	   notification setting
	   null default all notifications
	   1 email only
	   2 desktop only
	   3 disabled
	*/
	NotificationSetting NotificationSettingType `orm:"null"` //EMAIL_ONLY, DESKTOP_ONLY, NOTIFICATION_DISABLED

	/*
	   email notification triggered in queue
	*/
	NotificationTriggered time.Time `orm:"null"`

	//last access
	CountryCode string `orm:"size(2);null"` //TODO: country code choices limitation/validation
	City        string `orm:"size(100);null"`
	Device      string `orm:"size(30);null"` //TODO: device choices limitation/validation

	/* sender's email client */
	Client   EmailClientType `orm:"null"` //GMAIL_CLIENT, OUTLOOK_CLIENT
	Created  time.Time       `orm:"type(datetime);auto_now_add"`
	Modified time.Time       `orm:"type(datetime);auto_now"`
}
