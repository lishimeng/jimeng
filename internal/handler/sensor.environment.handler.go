package handler

import (
	"encoding/base64"
	log "github.com/jeanphorn/log4go"
	"github.com/lishimeng/jimeng/internal/etc"
	"github.com/lishimeng/jimeng/internal/influx"
	"github.com/lishimeng/jimeng/internal/mq"
	"github.com/lishimeng/jimeng/internal/tlv"
	"strconv"
)

func EnvironmentSensorHandler(session mq.Session, payload mq.PayloadRx) {

	log.Debug("environment sensor handler: %s", payload.Data)
	payloadTlv, err := base64.StdEncoding.DecodeString(payload.Data)
	if err != nil {
		log.Info(err)
		return
	}

	tags := map[string]string{
		"applicationID": payload.ApplicationID,
		"applicationName": payload.ApplicationName,
		"deviceName": payload.DeviceName,
		"devEUI": payload.DevEUI,
	}
	fields := make(map[string]interface{})

	frame := tlv.DecodeToFrame(payloadTlv[1:])

	for _, tag := range frame.Tags {

		switch tag.TagId {
		case 1:// speed
			fields["speed"] = tag.TagData
			break
		case 2:// light
			status := tag.TagData == "0"

			// 只处理 指定开关 为 有状态开关
			if payload.DevEUI == etc.Config.Mqtt.SwitchId {
				if status {
					SetCurrentState(true)
				} else {
					SetCurrentState(false)
				}
			}
			if status {
				fields["switch_state"] = 1
			} else {
				fields["switch_state"] = 0
			}
			log.Debug("switch status: %s:%b", payload.DeviceName, Shadow.CurrentState)
			break
		case 3:// message
			break
		case 5:// temporature
			val, err := strconv.ParseInt(tag.TagData, 10, 64)
			if err != nil {
				log.Info("%s\n", err)
			} else {
				value := val * 1600 / 4095
				temp := float64(value) / 10 - 40
				//f := fmt.Sprintf("%.1f", temp)
				fields["temperature"] = temp
			}
			break
		case 4:// humidity
			val, err := strconv.ParseInt(tag.TagData, 10, 64)
			if err != nil {
				log.Info("%s", err)
			} else {
				value := val * 1000 / 4095
				temp := float64(value) / 10
				//f := fmt.Sprintf("%.1f", temp)
				fields["humidity"] = temp
			}
			break
		case 6:// vibration
			val, err := strconv.ParseInt(tag.TagData, 10, 64)
			if err != nil {
				log.Info("%s", err)
			} else {
				value := float64(val) * 12.5 * 20 / 4096
				//f := fmt.Sprintf("%.1f", temp)
				fields["vibration"] = value
			}
			break
		case 7:// radar
			val, err := strconv.ParseInt(tag.TagData, 10, 64)
			if err != nil {
				log.Info("%s", err)
			} else {
				value := float64(val)
				log.Debug("radar: %d", value)
				fields["radar"] = value
				if val > 90 && val < 800 {
					go StatelessSyncState(session)// radar alarm
				}
			}
			break
		default:
			break
		}
	}

	if len(fields) > 0 {
		go influx.Save("environment_sensor", tags, fields)
	}
}