/*
  HTTP test refer to
  http://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest

  Websocket test refer to
  https://github.com/gorilla/websocket/blob/master/client_server_test.go
*/
package test

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	. "zenapi/models"
	// "github.com/dchest/uniuri"
	// "github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	_ "zenapi/routers"
)

func TestSignal(t *testing.T) {
	token := "wrghkjrgklkj4l2kl45"
	postData := url.Values{}
	postData.Set("sender", "sender@asdf.com")
	postData.Add("token", token)
	postData.Add("subject", "Hello Nice")
	postData.Add("body", "Hello Nice, this is James")
	postData.Add("to", `Ha <user1@asdf.com>,"M A"<user2@asdf.com>`)
	postData.Add("cc", "user3@asdf.com")
	postData.Add("bcc", "")
	postData.Add("links", `[{"url":"http://www.google.com", "urlHash":"-57483835", "urlDecode":"http://www.google.com", "plain": true}]`)
	postData.Add("tz_offset", "9")
	postData.Add("timezoneinfo", "Japan")
	beego.Debug(postData.Encode())
	req, _ := http.NewRequest("POST", "/zenapi/signal", bytes.NewBufferString(postData.Encode()))
	// req, _ := http.NewRequest("POST", "/zenapi/signal", strings.NewReader(postData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, req)

	beego.Trace("testing", "TestMain", "Code[%d]\n%s", w.Code, w.Body.String())
	fmt.Println(w.Body.String())

	Convey("Subject: Test Signal Endpoint\n", t, func() {
		Convey("Add Signal Return Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Add Signal Return Signal ID", func() {
			beego.Trace("testing", "TestMain", "Code[%d]\n%s", w.Code, w.Body.String())
			So(strings.Contains(w.Body.String(), "Id"), ShouldEqual, true)
			So(strings.Contains(w.Body.String(), token), ShouldEqual, true)
		})
		o := orm.NewOrm()
		Convey("Links Created", func() {
			var links []*Link
			qs1 := o.QueryTable("link")
			qs1.Filter("URL", "http://www.google.com").RelatedSel().All(&links)
			So(links[0].URL, ShouldEqual, "http://www.google.com")
			So(links[0].Signal.Id, ShouldEqual, 1)
			So(links[0].Signal.Token, ShouldEqual, token)
		})
		Convey("Signal Links", func() {
			s := new(Signal)
			qs2 := o.QueryTable(s)
			qs2.Filter("token", token).RelatedSel().One(s)
			o.LoadRelated(s, "Links")
			So(s.Links[0].URL, ShouldEqual, "http://www.google.com")
		})
	})
}
