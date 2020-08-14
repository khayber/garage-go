package main

import (
	"flag"
	"sync"
)

var gpioControlPin = flag.Int("control", 4, "GPIO Control Pin")
var gpioSensorPin = flag.Int("sensor", 17, "GPIO Sensor Pin")

var restEnable = flag.Bool("rest", false, "Enable REST API")
var restUser = flag.String("user", "god", "REST API Username")
var restPass = flag.String("pass", "damn", "REST API Password")
var restPort = flag.Int("port", 8080, "REST API port")
var restSsl = flag.Bool("ssl", false, "Enable SSL for REST API")

var telegramEnable = flag.Bool("telegram", false, "Enable Telegram Bot API")
var telegramToken = flag.String("token", "Your Token Here", "Telegram Bot API Token")
var telegramUser = flag.String("tg_user", "", "Telegram Username")

var monitorAutoclose = flag.Bool("autoclose", false, "Enable Auto Close feature")
var monitorClosetime = flag.Float64("closetime", 60, "Number of minutes after which door is closed")

var debug = flag.Bool("debug", false, "Enable debug logging")

// DEBUG controls extra logging
var DEBUG = false

func main() {
	flag.Parse()
	DEBUG = *debug
	var wg sync.WaitGroup

	door, _ := NewDoor(*gpioControlPin, *gpioSensorPin)
	if *monitorAutoclose {
		go door.monitor(*monitorClosetime)
		defer door.cleanup()
		wg.Add(1)
	}
	if *telegramEnable {
		bot, _ := NewDoorBot(door, *telegramToken, *telegramUser)
		go bot.Start()
		wg.Add(1) //make sure we don't exit if the rest server isn't configured
	}
	if *restEnable {
		rest, _ := NewRestService(door, *restUser, *restPass, *restPort, *restSsl)
		go rest.Listen()
		wg.Add(1)
	}

	wg.Wait()
}
