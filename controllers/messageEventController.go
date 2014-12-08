package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"

	"zenapi/models"
)

type Client struct {
	ClientID string
	Conn     *websocket.Conn
}

type MessageEventController struct {
	beego.Controller
}

// JoinClient method handles create and register websocket client

func (this *MessageEventController) JoinClient() {
	beego.Debug("JoinClient Received")
	clientID := this.GetString("client_id")
	if len(clientID) == 0 {
		beego.Error("No Client ID")
		http.Error(this.Ctx.ResponseWriter, "No Client ID", 404)
		return
	}

	//Upgrade from http to websocket
	//TODO: review "1024", read and write buffer size
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	// Join
	joinClient(clientID, ws)
	defer leaveClient(clientID)
	this.ServeJson()
}

// Get method receive events to deliver to client
func (this *MessageEventController) Post() {
	clientID := this.GetString("client_id")
	beego.Info(clientID)
	if len(clientID) == 0 {
		beego.Error("No Client ID", clientID)
		http.Error(this.Ctx.ResponseWriter, "No Client ID", 404)
		return
	}
	msg := this.GetString("msg")
	messageEventChannel <- newMessageEvent(models.EVENT_MESSAGE, clientID, msg)
	this.ServeJson()
}

func deliverMessageOverWebsocket(me models.MessageEvent) {
	data, err := json.Marshal(me)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	if sub, ok := joinedClients[me.ClientID]; ok {
		ws := sub.Conn
		if ws != nil {
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				//Client disconnected
				leavingChannel <- me.ClientID
			}
		} else {
			beego.Warning("Client Websocket Connection is nil:", me.ClientID)
		}
	} else {
		beego.Warning("Client ID not in joinedClients:", me.ClientID, len(joinedClients))
	}

	// for clientID, sub := range joinedClients {
	//  ws := sub.Conn
	//  if ws != nil {
	//    if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
	//      //User disconnected
	//      leavingChannel <- clientID
	//    }
	//  }
	// }
}

func newMessageEvent(et models.MessageEventType, clientID, msg string) models.MessageEvent {
	return models.MessageEvent{et, clientID, int(time.Now().Unix()), msg}
}

func joinClient(clientID string, ws *websocket.Conn) {
	joinedClientChannel <- Client{ClientID: clientID, Conn: ws}
}

func leaveClient(clientID string) {
	leavingChannel <- clientID
}

var (
	//Channel for new join client
	joinedClientChannel = make(chan Client, 1000)
	//Channel for exit user
	leavingChannel = make(chan string, 1000)
	//Channel to deliver message
	messageEventChannel = make(chan models.MessageEvent, 1000)
	//list of client who subscribes
	joinedClients = make(map[string]Client)
)

func messageOffice() {
	for {
		select {
		case client := <-joinedClientChannel:
			if !isClientExist(joinedClients, client.ClientID) {
				//New client
				joinedClients[client.ClientID] = client
				//we don't need to deliver this join event to everybody
				beego.Info("New client:", client.ClientID, ";Websocket:", client.Conn != nil)
			} else {
				//Old client
				beego.Info("Old client:", client.ClientID, ";Websocket:", client.Conn != nil)
			}
		case messageEvent := <-messageEventChannel:
			deliverMessageOverWebsocket(messageEvent)
			if messageEvent.Type == models.EVENT_MESSAGE {
				beego.Info("Message from:", messageEvent.ClientID, ";Message:", messageEvent.Message)
			}
		case leavingClientID := <-leavingChannel:
			if client, ok := joinedClients[leavingClientID]; ok {
				ws := client.Conn
				if ws != nil {
					ws.Close()
					//we don't need to deliver this leaving event to everybody
					beego.Info("Websocket closed:", leavingClientID)
				}
				delete(joinedClients, leavingClientID)
			}
		}
	}
}

func init() {
	go messageOffice()
}

func isClientExist(joinedClients map[string]Client, clientID string) bool {
	if _, ok := joinedClients[clientID]; ok {
		return true
	}
	return false
}
