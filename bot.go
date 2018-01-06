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

    // This button will be displayed in user's
    // reply keyboard.
    statusBtn := tb.ReplyButton{Text: "Status"}
    holdBtn := tb.ReplyButton{Text: "Hold"}
    openBtn := tb.ReplyButton{Text: "Open"}
    closeBtn := tb.ReplyButton{Text: "Close"}
    replyKeys := [][]tb.ReplyButton{
        []tb.ReplyButton{statusBtn},
        []tb.ReplyButton{holdBtn},
        []tb.ReplyButton{openBtn},
        []tb.ReplyButton{closeBtn},
    }

    // And this one — just under the message itself.
    // Pressing it will cause the client to send
    // the bot a callback.
    //
    // Make sure Unique stays unique as it has to be
    // for callback routing to work.
    inlineBtn := tb.InlineButton{
        Unique: "sad_moon",
        Text: "🌚 Button #2",
    }
    inlineKeys := [][]tb.InlineButton{
        []tb.InlineButton{inlineBtn},
        // ...
    }

    bot.Handle(&inlineBtn, func(c *tb.Callback) {
        log.Printf("callback %v", c)
        // on inline button pressed (callback!)

        // always respond!
        // b.Respond(c, &tb.CallbackResponse{...})
    })

    // Command: /start <PAYLOAD>
    bot.Handle("/start", func(m *tb.Message) {
        if !m.Private() {
            return
        }

        bot.Send(m.Sender, "Hello!", &tb.ReplyMarkup{
            ReplyKeyboard:  replyKeys,
            InlineKeyboard: inlineKeys,
        })
    })

    bot.Handle(&statusBtn, func(m *tb.Message) {
        log.Printf("messsage %v %v", m.Text, m.Sender)
        bot.Send(m.Sender, check_door() )
    })

    bot.Handle(&openBtn, func(m *tb.Message) {
        log.Printf("messsage %v %v", m.Text, m.Sender)
        bot.Send(m.Sender, open_door() )
    })

    bot.Handle(&closeBtn, func(m *tb.Message) {
        log.Printf("messsage %v %v", m.Text, m.Sender)
        bot.Send(m.Sender, close_door() )
    })

    bot.Handle(&holdBtn, func(m *tb.Message) {
        log.Printf("messsage %v %v", m.Text, m.Sender)
        bot.Send(m.Sender, hold_door() )
    })

    bot.Start()
}
