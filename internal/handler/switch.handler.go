package handler

import (
	"encoding/json"
	"github.com/lishimeng/jimeng/internal/etc"
	"github.com/lishimeng/jimeng/internal/mq"
	"time"
)

type SwitchShadow struct {

	CurrentState bool
	ChangeState bool
	PublishTimes int
}

// 下发次数
const MaxPublishTimes = 3
// 每次下发间隔
const PublishInterval = 1
// switch open 状态持续时间(下发操作执行完成后)
const AlarmDuring = 10

var firstState = true

var Shadow = SwitchShadow{
	CurrentState: false,
	ChangeState: false,
	PublishTimes: MaxPublishTimes,
}

func SetCurrentState(state bool) {

	Shadow.CurrentState = state
	if firstState {
		firstState = false
		Shadow.ChangeState = Shadow.CurrentState
	}
}

func SwitchStateSync() bool {
	return Shadow.ChangeState == Shadow.CurrentState
}

func SyncState(session mq.Session) {
	Shadow.PublishTimes = 0
	go _syncState(session)
}

func _syncState(session mq.Session) {

	for ;Shadow.PublishTimes < MaxPublishTimes; {

		if SwitchStateSync() {
			Shadow.PublishTimes = MaxPublishTimes
			return
		}
		Shadow.PublishTimes++
		// TODO wanted -> current
		tx := mq.PayloadTx{
			FPort: 2,
		}
		if Shadow.ChangeState {
			tx.Data = "MQIAAQA="
		} else {
			tx.Data = "MQIAAQE="
		}
		tmp, err := json.Marshal(tx)
		if err != nil {
			return
		}
		payload := string(tmp)
		session.Publish(etc.Config.Mqtt.SwitchAppId, etc.Config.Mqtt.SwitchId, payload)
		time.Sleep(time.Duration(4 * time.Second))
	}
}
