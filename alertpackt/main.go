// Package main provides ...
package main

import "log"
import "gopkg.in/telegram-bot-api.v4"

import "github.com/ewok/tools/packfreebook"

func main() {

	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	bot.Send(tgbotapi.NewMessage(174165027, "Started"))

	for update := range updates {
		reply := ""
		if update.Message == nil {
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, packfreebook.PackFreeBook())
		// println(packfreebook.PackFreeBook())

		switch update.Message.Command() {
		// удалить из списка
		case "start":
			reply = "Monitoring started"
		case "stop":
			reply = "Monitoring stopped"
		case "help":
			reply = "/start <time>;/stop"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, packfreebook.PackFreeBook())
		bot.Send(msg)
	}
	// for n := range devs {
	// 	if devs[n].Nickname != "" {
	// 		err = pb.PushNote(devs[n].Iden, "PacktFreeBook", packfreebook.PackFreeBook())
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }
}
