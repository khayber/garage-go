package main

import (
    "log"
    "time"
    tb "gopkg.in/tucnak/telebot.v2"
)

type MyBot struct {
    bot *tb.Bot
    user *tb.User
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
    } else {
        return iClosedKeys
    }
}

func initKeys() {
    iStatusBtn = tb.InlineButton{Unique: "status", Text: "❓ Status"}
    iOpenBtn = tb.InlineButton{Unique: "open", Text: "👍 Open"}
    iCloseBtn = tb.InlineButton{Unique: "close", Text: "👎 Close"}
    iHoldBtn = tb.InlineButton{Unique: "hold", Text: "✊ Hold"}
    iOpenKeys = [][]tb.InlineButton{
        []tb.InlineButton{iStatusBtn, iCloseBtn, iHoldBtn},
    }
    iClosedKeys = [][]tb.InlineButton{
        []tb.InlineButton{iStatusBtn, iOpenBtn},
    }
}

func NewMyBot(door *Door, token string, username string) (*MyBot, error) {
    mybot := &MyBot{}

    bot, err := tb.NewBot(tb.Settings{
                Token: token,
                Poller: &tb.LongPoller{10 * time.Second, -1}})
    if err != nil {
        log.Panic(err)
    }
    mybot.bot = bot

    log.Printf("Authorized on account %v", bot.Me)

    initKeys()

    bot.Handle(&iStatusBtn, func(c *tb.Callback) {
        log.Printf("Status %v", c.Sender)
        bot.Respond(c, &tb.CallbackResponse{c.ID, "", false, ""})
        msg, status := door.check();
        if status {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iOpenKeys})
        } else {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iClosedKeys})
        }
    })

    bot.Handle(&iOpenBtn, func(c *tb.Callback) {
        log.Printf("Open %v", c.Sender)
        bot.Respond(c, &tb.CallbackResponse{c.ID, "", false, ""})
        for msg := range door.open() {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iOpenKeys})
        }
    })

    bot.Handle(&iCloseBtn, func(c *tb.Callback) {
        log.Printf("Close %v", c.Sender)
        bot.Respond(c, &tb.CallbackResponse{c.ID, "", false, ""})
        for msg := range door.close() {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iClosedKeys})
        }
    })

    bot.Handle(&iHoldBtn, func(c *tb.Callback) {
        log.Printf("Hold %v", c.Sender)
        bot.Respond(c, &tb.CallbackResponse{c.ID, "", false, ""})
        msg, isOpen := door.hold()
        bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: getKeys(isOpen)})
    })

    bot.Handle("/start", func(m *tb.Message) {
        if !m.Private() {
            return
        }
        msg, isOpen := door.check()
        bot.Send(m.Sender, "Hello! " + msg, &tb.ReplyMarkup{InlineKeyboard: getKeys(isOpen)})
    })

    return mybot, nil
}

func (bot *MyBot) Start() {
    bot.bot.Start()
}

