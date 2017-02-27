package main

import (
    "log"
    "time"
    "github.com/tucnak/telebot"
)

func tgbot(token string) {
    bot, err := telebot.NewBot(token)
    if err != nil {
        log.Panic(err)
    }

    //bot.Debug = DEBUG
    log.Printf("Authorized on account %v", bot.Identity)

    messages := make(chan telebot.Message, 100)
    bot.Listen(messages, 60 * time.Second)

    for message := range messages {
        if message.Text == "" {
            continue
        }

        log.Printf("[%s] %s", message.Sender.Username, message.Text)

        switch cmd := message.Text; cmd {
            case "/start":
                bot.SendMessage(message.Chat, "Welcome", &telebot.SendOptions{
                    ReplyMarkup: telebot.ReplyMarkup{
                        ForceReply: true,
                        Selective: true,
                        CustomKeyboard: [][]string{
                            []string{"Open"},
                            []string{"Close"},
                            []string{"Status"},
                        },
                    },
                })
            case "Status": fallthrough
            case "/status":
                msg := check_door()
                bot.SendMessage(message.Chat, msg, nil)
            case "Open": fallthrough
            case "/open":
                msg := open_door()
                bot.SendMessage(message.Chat, msg, nil)
            case "Close": fallthrough
            case "/close":
                msg := close_door()
                bot.SendMessage(message.Chat, msg, nil)
            default:
                msg := "huh???"
                bot.SendMessage(message.Chat, msg, &telebot.SendOptions{
                    ReplyTo: message,
                })
        }
    }
}
