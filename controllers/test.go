/*
  HTTP test refer to
  http://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest

  Websocket test refer to
  https://github.com/gorilla/websocket/blob/master/client_server_test.go
*/
package controllers

import (
	// "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "os"
	"testing"
	. "zenapi/models"
	// "github.com/dchest/uniuri"
	// "github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"zenapi/zen"
)

func initTest() {
	zen.InitEnv()
	orm.RegisterModel(new(Access), new(Link), new(Receiver), new(Setting),
		new(Signal), new(User), new(UserAgent), new(UserTrack))
	zen.SyncDB()
}

func TestSignals(t *testing.T) {
	initTest()
	beego.Trace("TestCreateLink")
	signal := Signal{}
	Convey("Test create links", t, func() {
		// linksJsonString := `[{'url': 'http://www.abc.com', 'urlHash': '-16913783', 'urlDecoded': 'http://www.abc.com', 'plain': false}]`
		var linksJsonString = `[{"url": "http://www.abc.com", "urlHash": "-16913783", "urlDecoded": "http://www.abc.com", "plain": false}]`
		links := createLinks(linksJsonString, &signal)
		So(links[0].URL, ShouldEqual, "http://www.abc.com")
		So(links[0].URLHash, ShouldEqual, "-16913783")
		So(links[0].Signal, ShouldEqual, &signal)
	})
	// fmt.Println("TestSaveAndAssignReceivers")
	Convey("TestSaveAndAssignReceivers", t, func() {
		receiverEmailStrings := make(map[string]string)
		receiverEmailStrings["to"] = `MC HotDog <c1@asdf.com>,c2@asdf.com`
		receiverEmailStrings["cc"] = `JJ <jj@asdf.com>`
		receiverEmailStrings["bcc"] = `jk@asdf.com`
		saveAndAssignReceivers(receiverEmailStrings, &signal)
		So(signal.ReceiverEmails, ShouldEqual, "c1@asdf.com;c2@asdf.com;jj@asdf.com;jk@asdf.com")
	})
}
