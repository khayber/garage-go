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

    iStatusBtn := tb.InlineButton{Unique: "status", Text: "❓ Status"}
    iOpenBtn := tb.InlineButton{Unique: "open", Text: "👍 Open"}
    iCloseBtn := tb.InlineButton{Unique: "close", Text: "👎 Close"}
    iHoldBtn := tb.InlineButton{Unique: "hold", Text: "✊ Hold"}
    iOpenKeys := [][]tb.InlineButton{
        []tb.InlineButton{iStatusBtn, iCloseBtn, iHoldBtn},
    }
    iClosedKeys := [][]tb.InlineButton{
        []tb.InlineButton{iStatusBtn, iOpenBtn},
    }

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
        msg, status := door.hold()
        if status {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iOpenKeys})
        } else {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iClosedKeys})
        }
    })

    bot.Handle("/start", func(m *tb.Message) {
        if !m.Private() {
            return
        }

        msg, status := door.check()

        bot.Send(m.Sender, "Hello! " + msg, &tb.ReplyMarkup{
            InlineKeyboard: func() [][]tb.InlineButton {
                if status {return iOpenKeys} else { return iClosedKeys }
            }(),
        })
    })

    return mybot, nil
}


func (bot *MyBot) Start() {
    bot.bot.Start()
}

