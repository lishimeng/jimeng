package handler

import (
	"encoding/json"
	"github.com/lishimeng/jimeng/internal/mq"
)


type MessageHandler func (session mq.Session, payload mq.PayloadRx)

var Handlers = make(map[string]MessageHandler)

func Callback(session mq.Session, topic string, msg []byte)() {

	var payload mq.PayloadRx
	_ = json.Unmarshal(msg, &payload)

	appId := payload.ApplicationID
	if v, ok := Handlers[appId]; ok {
		v(session, payload)
	}
}