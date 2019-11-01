package handler

import (
	"encoding/json"
	log "github.com/jeanphorn/log4go"
	"github.com/lishimeng/jimeng/internal/etc"
	"github.com/lishimeng/jimeng/internal/mq"
	"time"
)

var radarAlarmLock = false

func StatelessSyncState(session mq.Session) {

	if radarAlarmLock {
		return
	}
	log.Debug("radar alarm begin..........")
	radarAlarmLock = true
	_statelessChangeState(session, "MQIAAQA=")
	time.Sleep(time.Duration(AlarmDuring * time.Second))
	_statelessChangeState(session, "MQIAAQE=")
	radarAlarmLock = false
	log.Debug("radar alarm stop..........")
}

func _statelessChangeState(session mq.Session, state string) {

	publishTimes := 0
	tx := mq.PayloadTx{
		FPort: 2,
	}

	for ;publishTimes < MaxPublishTimes; {
		publishTimes++
		tx.Data = state
		tmp, err := json.Marshal(tx)
		if err != nil {
			return
		}
		payload := string(tmp)
		session.Publish(etc.Config.Mqtt.SwitchAppId, etc.Config.Mqtt.SwitchId, payload)
		time.Sleep(time.Duration(PublishInterval * time.Second))
	}
}