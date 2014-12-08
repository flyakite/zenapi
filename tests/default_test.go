/*
	HTTP test refer to
	http://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest

	Websocket test refer to
	https://github.com/gorilla/websocket/blob/master/client_server_test.go
*/
package test

import (
	"bytes"
	"github.com/astaxie/beego"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	// "zenapi/controllers"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	_ "zenapi/routers"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestMain is a sample to run an endpoint test
func TestMain(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestMain", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestMessageEventClientNotJoined(t *testing.T) {
	//http://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest
	//"client_id": {"asdfasdf"}, "msg": {"msg"}
	data := url.Values{}
	data.Set("client_id", "clientID_001")
	data.Add("msg", "test")
	r, _ := http.NewRequest("POST", "/messageevent", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestMessageEventClientNotJoined", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

type Server struct {
	*httptest.Server
	URL string
}

var cstDialer = websocket.Dialer{
	Subprotocols:    []string{"p1", "p2"},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type cstHandler struct{ *testing.T }

func (t cstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	beego.Debug("ServeHTTP", r.URL)
	beego.BeeApp.Handlers.ServeHTTP(w, r)
}

func newServer(t *testing.T) *Server {
	var s Server
	s.Server = httptest.NewServer(cstHandler{t})
	s.URL = "ws" + s.Server.URL[len("http"):]
	return &s
}

func TestMessageEventClientJoined(t *testing.T) {
	s := newServer(t)
	defer s.Close()
	beego.Debug(s.URL)
	ws, _, err := cstDialer.Dial(s.URL+"/messageevent/joinclient?client_id=clientID_001", nil)
	defer ws.Close()
	if err != nil {
		t.Fatal("Dial: %v", err)
	}

	data := url.Values{}
	data.Set("client_id", "clientID_001")
	data.Add("msg", "test")
	r, _ := http.NewRequest("POST", s.URL+"/messageevent", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	beego.Trace("testing", "TestMessageEventClientJoined", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Joined Websocket Client Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})

}
