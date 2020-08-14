package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// DoorBot defines the bot for the door
type DoorBot struct {
	bot     *tb.Bot
	user    *tb.User
	message *tb.Message
}

var iStatusBtn tb.InlineButton
var iOpenBtn tb.InlineButton
var iCloseBtn tb.InlineButton
var iHoldBtn tb.InlineButton
var iOpenKeys [][]tb.InlineButton
var iClosedKeys [][]tb.InlineButton

func getKeys(isOpen bool) [][]tb.InlineButton {
	if isOpen {
		return iOpenKeys
	}
	return iClosedKeys
}

func initKeys() {
	iStatusBtn = tb.InlineButton{Unique: "status", Text: "❓ Status"}
	iOpenBtn = tb.InlineButton{Unique: "open", Text: "👍 Open"}
	iCloseBtn = tb.InlineButton{Unique: "close", Text: "👎 Close"}
	iHoldBtn = tb.InlineButton{Unique: "hold", Text: "✊ Hold"}
	iOpenKeys = [][]tb.InlineButton{
		{iStatusBtn, iCloseBtn, iHoldBtn},
	}
	iClosedKeys = [][]tb.InlineButton{
		{iStatusBtn, iOpenBtn},
	}
}

// NewDoorBot creates and initialie a new DoorBot
func NewDoorBot(door *Door, token string, username string) (*DoorBot, error) {
	mybot := &DoorBot{}

	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second, LastUpdateID: -1}})
	if err != nil {
		log.Panic(err)
	}
	mybot.bot = bot

	log.Printf("Authorized on account %v", bot.Me)

	initKeys()

	bot.Handle(&iStatusBtn, func(c *tb.Callback) {
		log.Printf("Status %v", c.Sender)
		if c.Sender.Username != username {
			log.Printf("I've been HACKED!!!")
			return
		}
		bot.Respond(c, &tb.CallbackResponse{
			CallbackID: c.ID,
			Text:       "",
			ShowAlert:  false,
			URL:        ""})
		msg, status := door.check()
		if status {
			bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iOpenKeys})
		} else {
			bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iClosedKeys})
		}
	})

	bot.Handle(&iOpenBtn, func(c *tb.Callback) {
		log.Printf("Open %v", c.Sender)
		if c.Sender.Username != username {
			log.Printf("I've been HACKED!!!")
			return
		}
		bot.Respond(c, &tb.CallbackResponse{
			CallbackID: c.ID,
			Text:       "",
			ShowAlert:  false,
			URL:        ""})
		for msg := range door.open() {
			bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iOpenKeys})
		}
	})

	bot.Handle(&iCloseBtn, func(c *tb.Callback) {
		log.Printf("Close %v", c.Sender)
		if c.Sender.Username != username {
			log.Printf("I've been HACKED!!!")
			return
		}
		bot.Respond(c, &tb.CallbackResponse{
			CallbackID: c.ID,
			Text:       "",
			ShowAlert:  false,
			URL:        ""})
		for msg := range door.close() {
			bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iClosedKeys})
		}
	})

	bot.Handle(&iHoldBtn, func(c *tb.Callback) {
		log.Printf("Hold %v", c.Sender)
		if c.Sender.Username != username {
			log.Printf("I've been HACKED!!!")
			return
		}
		bot.Respond(c, &tb.CallbackResponse{
			CallbackID: c.ID,
			Text:       "",
			ShowAlert:  false,
			URL:        ""})
		msg, isOpen := door.hold()
		bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: getKeys(isOpen)})
	})

	bot.Handle("/start", func(m *tb.Message) {
		log.Printf("Hold %v", m.Sender)
		if m.Sender.Username != username {
			log.Printf("I've been HACKED!!!")
			return
		}
		if !m.Private() {
			return
		}
		msg, isOpen := door.check()
		bot.Send(m.Sender, "Hello! "+msg, &tb.ReplyMarkup{InlineKeyboard: getKeys(isOpen)})
	})

	return mybot, nil
}

// Start is used to start stuff...
func (bot *DoorBot) Start() {
	bot.bot.Start()
}
