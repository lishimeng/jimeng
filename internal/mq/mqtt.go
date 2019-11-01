package mq

import "github.com/eclipse/paho.mqtt.golang"

type MessageCallback func(session Session, topic string, msg []byte)

type Config struct {
	Broker string
	ClientId string
	Topics string
	Publish string
	DbUrl string
	DbName string
}

type Session struct {
	client *mqtt.Client
	Config *Config
	OnMessage MessageCallback
}

type PayloadRx struct {
	ApplicationID string `json:"applicationID"`
	ApplicationName string `json:"applicationName"`
	DeviceName string `json:"deviceName"`
	DevEUI string `json:"devEUI"`
	Data string `json:"data"`
}

type PayloadTx struct {
	FPort int32 `json:"fPort"`
	Data string `json:"data"`
}

type Handler interface {
	/* 上传 */
	Publish(device string, msg string)
	/* 订阅 */
	Subscribe(callback MessageCallback)
	/* 取消订阅 */
	Unsubscribe()
	/* 连接 */
	Connect()
	/* 断开 */
	Close()
	/* 消息回调 */
	MessageHandler(client mqtt.Client, msg mqtt.Message)
}