package main

import (
    "os"
    "log"
    "gopkg.in/telegram-bot-api.v4"
)

func tgbot() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = false

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        switch cmd := update.Message.Text; cmd {
            case "/start":
                //TODO - create a custom keyboard
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome.")
                bot.Send(msg)
            case "/check": fallthrough
            case "/status":
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, check_door())
                bot.Send(msg)
            case "/open":
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, open_door())
                bot.Send(msg)
            case "/close":
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, close_door())
                bot.Send(msg)
            default:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "???")
                msg.ReplyToMessageID = update.Message.MessageID
                bot.Send(msg)
        }
    }
}
