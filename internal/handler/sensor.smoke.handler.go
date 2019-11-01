package handler

import (
	"encoding/base64"
	log "github.com/jeanphorn/log4go"
	"github.com/lishimeng/jimeng/internal/mq"
)

func SmokeSensorHandler(session mq.Session, payload mq.PayloadRx) {

	log.Debug("smoke sensor handler: %s", payload.Data)
	payloadData, err := base64.StdEncoding.DecodeString(payload.Data)
	if err != nil {
		log.Info(err)
		return
	}

	if len(payloadData) == 0x12 {
		high := payloadData[8]
		low := payloadData[9]
		humi := (int16(high & 0xff) << 8) + int16(low & 0xff)
		humidity := float64(humi) / 10

		log.Debug("humidity:%f", humidity)
		if humidity > 30.0 {
			// too high
			//downControll(session, false)
		} else if humidity < 5.0 {
			// too low
			//downControll(session, true)
		} else {
			// Normal
		}
	}
}