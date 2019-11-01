package main

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/lishimeng/jimeng/internal/api"
	"github.com/lishimeng/jimeng/internal/etc"
	"github.com/lishimeng/jimeng/internal/handler"
	"github.com/lishimeng/jimeng/internal/mq"
	"time"
)

func main() {

	log.LoadConfiguration("log.json")
	defer log.Close()
	etc.SetConfigName("conf.toml")
	etc.AddEnvPath(".")
	etc.AddEnvPath("/etc/server")
	err := etc.LoadEnvs()
	if err != nil {
		log.Info("%s", err)
		return
	}

	configuration := etc.Config
	config := mq.Config{
		Broker: configuration.Mqtt.Broker,
		ClientId: configuration.Mqtt.ClientId,
		Topics: configuration.Mqtt.Subscribe,
		Publish: configuration.Mqtt.Upstream,
		DbName: configuration.Db.Name,
		DbUrl: configuration.Db.Url,
	}

	handler.Handlers[configuration.Mqtt.SwitchAppId] = handler.EnvironmentSensorHandler
	handler.Handlers[configuration.Mqtt.SoilAppId] = handler.SoilSensorHandler
	handler.Handlers[configuration.Mqtt.SmokeAppId] = handler.SmokeSensorHandler

	log.Info("Start MQTT handler")
	session := mq.CreateSession(&config, handler.Callback)
	log.Info("Connect:%s", config.Broker)
	session.Connect()
	log.Info("Subscribe the device messages")
	session.Subscribe(fmt.Sprintf(configuration.Mqtt.Subscribe, configuration.Mqtt.SoilAppId))
	session.Subscribe(fmt.Sprintf(configuration.Mqtt.Subscribe, configuration.Mqtt.SwitchAppId))
	log.Info("MQTT start successfully")
	log.Info("Start Web Server:%s", time.Now())

	webOptions := api.WebOptions{}
	webOptions.Listen = configuration.Web.Listen
	webEngine := api.Create(&webOptions)
	webEngine.Start()

}