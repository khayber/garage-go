package main

import (
    "flag"
    "sync"
)

var gpio_control_pin = flag.Int("control", 4, "GPIO Control Pin")
var gpio_sensor_pin = flag.Int("sensor", 17, "GPIO Sensor Pin")

var rest_enable = flag.Bool("rest", false, "Enable REST API")
var rest_user = flag.String("user", "god", "REST API Username")
var rest_pass = flag.String("pass", "damn", "REST API Password")
var rest_port = flag.Int("port", 8080, "REST API port")
var rest_ssl = flag.Bool("ssl", false, "Enable SSL for REST API")

var telegram_enable = flag.Bool("telegram", false, "Enable Telegram Bot API")
var telegram_token = flag.String("token", "Your Token Here", "Telegram Bot API Token")
var telegram_user = flag.String("tg_user", "", "Telegram Username")

var monitor_autoclose = flag.Bool("autoclose", false, "Enable Auto Close feature")
var monitor_closetime = flag.Float64("closetime", 60, "Number of minutes after which door is closed")

var debug = flag.Bool("debug", false, "Enable debug logging")
var DEBUG = false

func main() {
    flag.Parse()
    DEBUG = *debug
    var wg sync.WaitGroup

    door, _ := NewDoor(*gpio_control_pin, *gpio_sensor_pin)
    if *monitor_autoclose {
        go door.monitor(*monitor_closetime)
        wg.Add(1)
    }
    if *telegram_enable {
        bot, _ := NewMyBot(door, *telegram_token, *telegram_user)
        go bot.Start()
        wg.Add(1) //make sure we don't exit if the rest server isn't configured
    }
    if *rest_enable {
        rest, _ := NewRestService(door, *rest_user, *rest_pass, *rest_port, *rest_ssl)
        go rest.Listen()
        wg.Add(1)
    }

    wg.Wait()
}
