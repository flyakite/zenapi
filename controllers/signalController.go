package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	// "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strconv"
	"strings"
	. "zenapi/models"
	"zenapi/utils"
)

type SignalController struct {
	BaseController
}

func (this *SignalController) AddSignal() {
	beego.Info(this.Input())
	beego.Info(this.Ctx.Input.RequestBody)

	//client debug info
	version := this.GetString("version")
	beego.Info("client version: ", version)
	debuginfo := this.GetString("debuginfo")
	if debuginfo != "" {
		for debuginfoline := range strings.Split(debuginfo, "#_#") {
			beego.Error(debuginfoline)
		}
	}

	var signal Signal

	tzOffset := this.GetString("tz_offset")
	if tzOffset != "" {
		tzOffset64, err := strconv.ParseInt(tzOffset, 10, 0)
		if err != nil {
			beego.Error("TzOffsetError: ", err)
		} else {
			signal.TzOffset = int(tzOffset64) * 60
		}
	}
	timezoneinfo := this.GetString("timezoneinfo")
	if timezoneinfo != "" {
		beego.Info(timezoneinfo)
	}

	token := this.GetString("token")
	beego.Info("token: ", token)
	if token != "" {
		signal.Token = token
	} else {
		beego.Error("TokenInvalid: ", token)
		this.Data["json"] = `{"error":1}`
		this.ServeJson()
		return
	}
	sender := this.GetString("sender")
	if utils.IsValidEmail(sender) {
		beego.Info("Sender: ", sender)
		signal.Email = sender
	} else {
		beego.Error("SenderEmailNotValid: ", sender)
		this.Data["json"] = `{"error":1}`
		this.ServeJson()
		return
	}

	subject := this.GetString("subject")

	body := this.GetString("body")
	outlook := this.GetString("outlook")
	if outlook != "" {
		signal.Client = OUTLOOK_CLIENT
	}

	signal.Subject = subject
	signal.Body = body

	//receivers
	var receivers = make(map[string]string)
	for _, v := range SIGNAL_RECEIVER_KEYS {
		receivers[v] = this.GetString(v)
	}
	saveAndAssignReceivers(receivers, &signal)
	var db = orm.NewOrm()
	_, err := db.Insert(&signal)
	if err != nil {
		beego.Error("SignalInsertError: ", err)
		this.Data["json"] = `{"error":1}`
		this.ServeJson()
		return
	}
	linksJsonString := this.GetString("links")
	if linksJsonString != "" {
		beego.Info("Create links")
		createLinks(linksJsonString, &signal)
	}
	beego.Info("signal id: ", signal.Id)
	//var data map[string]interface{}
	//data["signal"] = make(map[string]interface{})
	this.Data["json"] = &signal
	this.ServeJson()
	// signalByte, err := json.Marshal(signal)
	// if err != nil {
	// }
	//TODO: update statitics
}

type LinkContainer struct {
	URL       string `json:"url"`
	URLHash   string `json:"urlHash"`
	URLDecode string `json:"urlDecode"`
	Plain     bool   `json:"plain"`
}

func createLinks(linksJsonString string, signal *Signal) []Link {
	beego.Info("createLinks")
	//[{u'url': u'http://www.abc.com', u'urlHash': u'-16913783', u'urlDecoded': u'http://www.abc.com', u'plain': False}]
	var LinkContainers = []LinkContainer{}
	json.Unmarshal([]byte(linksJsonString), &LinkContainers)
	var links []Link
	if len(LinkContainers) > 0 {
		for _, l := range LinkContainers {
			link := Link{}
			link.URL = l.URL
			link.URLHash = l.URLHash
			link.Signal = signal
			links = append(links, link)
		}
		var db = orm.NewOrm()
		_, err := db.InsertMulti(100, links)
		if err != nil {
			beego.Error(err)
			_, err := db.InsertMulti(1, links)
			if err != nil {
				beego.Error(err)
			}
		}
	} else {
		beego.Error("LinkUnmarshalError ", linksJsonString)
	}
	for _, l := range links {
		beego.Info("link id:", l.Id)
		beego.Info("link url:", l.URL)
		beego.Info("link signal:", l.Signal)
	}
	// o := orm.NewOrm()
	// qs := o.QueryTable(new(Link))
	// qs.Filter("URL")
	return links
}

func seperatefullEmailAddress(fullEmailAddress string) (string, string) {
	// 1. ld@ab.com
	// 2. ld@ab.com,
	// 3. Leo Del <ld@ab.com>
	// 4. "Leo Del" <ld@ab.com>
	// 5. "Leo%2C Del" <ld@ab.com>

	if fullEmailAddress == "" {
		return "", ""
	}

	r, _ := regexp.Compile("\"?(.*?)\"?\\s*\\<(.*)\\>")
	fullEmailAddress = strings.Trim(fullEmailAddress, ", ")
	if strings.HasSuffix(fullEmailAddress, ">") {
		m := r.FindStringSubmatch(fullEmailAddress)
		if len(m) == 3 {
			email, name := m[2], m[1]
			if utils.IsValidEmail(email) {
				return email, name
			} else {
				beego.Error("InvalidEmailString1 ", fullEmailAddress)
				return email, name //TODO: return empty string
			}
		} else {
			beego.Error("InvalidEmailString2 ", fullEmailAddress)
			return "", ""
		}
	} else {
		//pure email
		if utils.IsValidEmail(fullEmailAddress) {
			return fullEmailAddress, ""
		} else {
			beego.Error("InvalidEmailString3 ", fullEmailAddress)
			return "", ""
		}
	}
}

func saveAndAssignReceivers(receiverEmailStrings map[string]string, signal *Signal) {

	receiverArchive, _ := json.Marshal(receiverEmailStrings)
	signal.ReceiverArchive = string(receiverArchive)
	//if there's no ',' in user name
	// var person Receiver
	// var tempReceivers []*tempReceiver
	// var receiversMap = make(map[string]*Receiver)
	// var receiversByField = make(map[string][]*Receiver)
	var emails []string
	// var emailsByField = make(map[string][]string)
	// var email, name string
	// var receivers []Receiver

	for _, v := range receiverEmailStrings {
		emailList := strings.Split(v, ",")
		for _, fullEmailAddress := range emailList {
			email, _ := seperatefullEmailAddress(fullEmailAddress)

			// person = Receiver{email, name}
			// beego.Debug(person.Email, " ", person.Names)
			// receiversMap[email] = &person
			// receiversByField[k] = append(receiversByField[k], &person)
			//tempReceivers = append(receivers, &person)
			if email != "" {
				emails = append(emails, email)
			}
			//emailsByField[k] = append(emailsByField[k], email)
		}
		//Receiver().Filter("emails__in", emailsByField[k]...).All(&receivers)
	}
	signal.ReceiverEmails = strings.Join(emails, SIGNAL_RECEIVER_EMAILS_SEPERATOR)

	// qs := DBO.QueryTable(new(Receiver))
	// qs.Filter("email__in", emails...)
	// var receivers []Receiver
	// Receiver().Filter("emails__in", emails...).All(&receivers)
	// if len(receivers) > 0 {
	// 	var emailsCreated []string
	// 	for _, r := range receivers {
	// 		emailsCreated = append(emailsCreated, r.Email)
	// 	}
	// 	var receiversToCreate []Receiver
	// 	for _, e := range emails {
	// 		if !utils.StringinSlice(e, emailsCreated) {
	//        receiversToCreate = append(receiversToCreate, receiversMap[e])
	// 		}
	// 	}
	//    if len(receiversToCreate) > 0 {
	//      successNums, err := Receiver().InsertMulti(10, receiversToCreate)
	//      if err != nil {
	//        beego.Error(err)
	//        successNums, err := Receiver().InsertMulti(1, receiversToCreate)
	//      }
	//    }
	// }
	//  for k := range receiverEmailStrings{

	//  }
}

type SignalImageController struct {
	BaseController
}

func (this *SignalImageController) GetSignalImage() {
	token := this.GetString("t")
	//TODO: use cache to prevent hight frequent access
	if token == "" {
		this.ServePixelImage()
		return
	}
	beego.Info(token)
	signal := Signal{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Signal))
	err := qs.Filter("token", token).One(&signal)
	if err == orm.ErrMultiRows {
		qs.Filter("token", token).OrderBy("-created").Limit(1).One(&signal)
	} else if err == orm.ErrNoRows {
		beego.Error("SignalWithTokenNotFound: ", token)
		this.ServePixelImage()
		return
	}
	beego.Info(this.Ctx.Input.IP())
	beego.Info(this.Ctx.Input.UserAgent())
	beego.Info(this.Ctx.Input.Referer())

}

func recognizeUser(signal *Signal) {

	receiverEmails := strings.Split(signal.ReceiverEmails, SIGNAL_RECEIVER_EMAILS_SEPERATOR)
	beego.Info(len(receiverEmails), " ", receiverEmails)
	if len(receiverEmails) == 1 {
		//TODO
		return
	} else {
		return
	}
}

type ProxySignalController struct {
	BaseController
}

func (this *ProxySignalController) AddSignal() {
	beego.Info(this.Input())
	beego.Info(this.Ctx.Input.RequestBody)

	version := this.GetString("version")
	debuginfo := this.GetString("debuginfo")
	tzOffset := this.GetString("tz_offset")
	timezoneinfo := this.GetString("timezoneinfo")
	token := this.GetString("token")
	sender := this.GetString("sender")
	subject := this.GetString("subject")
	// body := this.GetString("body")
	outlook := this.GetString("outlook")
	links := this.GetString("links")
	debug := this.GetString("debug")
	var receivers = make(map[string]string)
	for _, v := range SIGNAL_RECEIVER_KEYS {
		receivers[v] = this.GetString(v)
	}

	trackerURL := "https://zenblip.appspot.com"
	path := "/signals/add"
	postData := url.Values{}
	postData.Set("version", version)
	postData.Add("debuginfo", debuginfo)
	postData.Add("tz_offset", tzOffset)
	postData.Add("timezoneinfo", timezoneinfo)
	postData.Add("token", token)
	postData.Add("sender", sender)
	postData.Add("subject", subject)
	postData.Add("outlook", outlook)
	postData.Add("links", links)
	if debug != "" {
		path = "/signals/debug"
	}
	for k, v := range receivers {
		postData.Add(k, v)
	}

	req, _ := http.NewRequest("POST", trackerURL+path, bytes.NewBufferString(postData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))

	client := &http.Client{}
	resp, _ := client.Do(req)
	beego.Info(resp.Status)
	this.Data["json"] = `{"success":1}`
	this.ServeJson()
	return
}
