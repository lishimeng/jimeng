package mq

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/jeanphorn/log4go"
	"os"
)


func CreateSession(config *Config, callback MessageCallback) Session {
	session := Session{}
	opts := MQTT.NewClientOptions().AddBroker(config.Broker)
	opts.SetClientID(config.ClientId)
	opts.SetDefaultPublishHandler(session.MessageHandler)

	c := MQTT.NewClient(opts)
	session.client = &c
	session.Config = config
	session.OnMessage = callback
	return session
}

func (session *Session) Publish(app string, device string, payload string) () {

	client := *session.client
	topic := fmt.Sprintf(session.Config.Publish, app, device)

	log.Info("publish msg: %s:%s", topic, payload)
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
}

func (session *Session) Subscribe(topic string) () {
	client := *session.client

	log.Info("Subscribe: %s", topic)
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Info(token.Error())
		os.Exit(1)
	}
}

func (session *Session) Unsubscribe() () {
	client := *session.client
	if token := client.Unsubscribe(session.Config.Topics); token.Wait() && token.Error() != nil {
		log.Info(token.Error())
		os.Exit(1)
	}
}

func (session *Session) Connect() () {
	client := *session.client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (session *Session) Close() () {
	client := *session.client
	client.Disconnect(250)
}

func (session *Session) MessageHandler(client MQTT.Client, msg MQTT.Message) {
	log.Fine("TOPIC: %s", msg.Topic())
	log.Fine("MSG RAW: %s", msg.Payload())
	session.OnMessage(*session, msg.Topic(), msg.Payload())
}
