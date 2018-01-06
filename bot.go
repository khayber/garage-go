package main

import (
    "log"
    "time"
    tb "gopkg.in/tucnak/telebot.v2"
)

func tgbot(token string, username string) {
    bot, err := tb.NewBot(tb.Settings{
                Token: token,
                Poller: &tb.LongPoller{10 * time.Second, -1}})
    if err != nil {
        log.Panic(err)
    }

    log.Printf("Authorized on account %v", bot.Me)

    iStatusBtn := tb.InlineButton{Unique: "status", Text: "❓Status"}
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
        msg, status := check_door();
        if status {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iOpenKeys})
        } else {
            bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iClosedKeys})
        }
    })

    bot.Handle(&iOpenBtn, func(c *tb.Callback) {
        log.Printf("Open %v", c.Sender)
        bot.Respond(c, &tb.CallbackResponse{c.ID, "", false, ""})
        msg, _ := open_door()
        bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iOpenKeys})
    })

    bot.Handle(&iCloseBtn, func(c *tb.Callback) {
        log.Printf("Close %v", c.Sender)
        bot.Respond(c, &tb.CallbackResponse{c.ID, "", false, ""})
        msg, _ := close_door()
        bot.Edit(c.Message, msg, &tb.ReplyMarkup{InlineKeyboard: iClosedKeys})
    })

    bot.Handle(&iHoldBtn, func(c *tb.Callback) {
        log.Printf("Hold %v", c.Sender)
        bot.Respond(c, &tb.CallbackResponse{c.ID, "", false, ""})
        msg, status := hold_door()
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

        msg, status := check_door()

        bot.Send(m.Sender, "Hello! " + msg, &tb.ReplyMarkup{
            InlineKeyboard: func() [][]tb.InlineButton {
                if status {return iOpenKeys} else { return iClosedKeys }
            }(),
        })
    })

    bot.Start()
}
