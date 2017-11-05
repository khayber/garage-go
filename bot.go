package main

import (
    "log"
    "time"
    "github.com/tucnak/telebot"
)

func tgbot(token string, username string) {
    bot, err := telebot.NewBot(token)
    if err != nil {
        log.Panic(err)
    }

    log.Printf("Authorized on account %v", bot.Identity)

    messages := make(chan telebot.Message, 100)
    bot.Listen(messages, 60 * time.Second)

    for message := range messages {
        if DEBUG {
            log.Printf("messsage %v", message)
        }

        log.Printf("[%s] %s", message.Sender.Username, message.Text)
        if message.Sender.Username != username {
            log.Printf("ERROR invalid user")
            continue;
        }
        switch cmd := message.Text; cmd {
            case "/start":
                bot.SendMessage(message.Chat, "Welcome", &telebot.SendOptions{
                    ReplyMarkup: telebot.ReplyMarkup{
                        ForceReply: true,
                        Selective: true,
                        CustomKeyboard: [][]string{
                            []string{"/status"},
                            []string{"/open"},
                            []string{"/close"},
                            []string{"/hold"},
                        },
                    },
                })
            case "/status":
                msg := check_door()
                bot.SendMessage(message.Chat, msg, nil)
            case "/open":
                msg := open_door()
                bot.SendMessage(message.Chat, msg, nil)
            case "/close":
                msg := close_door()
                bot.SendMessage(message.Chat, msg, nil)
            case "/hold":
                msg := hold_door()
                bot.SendMessage(message.Chat, msg, nil)
            default:
                msg := "huh???"
                bot.SendMessage(message.Chat, msg, &telebot.SendOptions{
                    ReplyTo: message,
                })
        }
    }
}
