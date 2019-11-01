package influx

import (
	"github.com/influxdata/influxdb1-client"
	log "github.com/jeanphorn/log4go"
	"github.com/lishimeng/jimeng/internal/etc"
	"net/url"
	"time"
)

var failTimes = 0

func Save(measurement string, tags map[string]string, fields map[string]interface{}) {

	host, err := url.Parse(etc.Config.Db.Url)
	if err != nil {
		log.Info(err)
		return
	}

	con, err := client.NewClient(client.Config{URL: *host})
	if err != nil {
		log.Info(err)
		return
	}

	point := client.Point{
		Measurement:measurement,
		Tags: tags,
		Fields: fields,
		Time: time.Now(),
	}

	points := []client.Point{point}

	bps := client.BatchPoints{
		Points: points,
		Database: etc.Config.Db.Name,
	}

	_, err = con.Write(bps)

	if err != nil {
		log.Info(err)
		failTimes++
	} else {
		failTimes = 0
	}

	if failTimes > 0x10 {
		panic("influx failed too much times")
	}
}