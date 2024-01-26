package telegram

import (
	"context"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	subscribeCommand = "subscrive"
	unsubscribeCommand = "unsubscrive"
)

type repository interface {
	AddUser(ctx context.Context, ID int64) error
	RemoveUser(ctx context.Context, ID int64) error
}
type TgBot struct {
	repo repository

	*tgbotapi.BotAPI
}

func (tg *TgBot) Serve(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{subscribeCommand, unsubscribeCommand}
	updates := tg.GetUpdatesChan(u)

	for {
		select {
			case update := <-updates:
				tg.processMessage(ctx, *update.Message)
			case <-ctx.Done():
				return
		}
	}
}

func (tg *TgBot) processMessage(ctx context.Context, msg tgbotapi.Message) {
	var err error

	switch msg.Command() {
		case subscribeCommand:
			err = tg.processSubscribeCommand(ctx, msg)
		case unsubscribeCommand:
			err = tg.processUnsubscribeCommand(ctx, msg)
		default:
			m := tgbotapi.NewMessage(msg.Chat.ID, "No such command.")
			tg.BotAPI.Send(m)
	}

	if err != nil {
		errMsg := tgbotapi.NewMessage(msg.Chat.ID, "Something went wrong. ٩(͡๏̯͡๏)۶")
		tg.BotAPI.Send(errMsg)
	}
}

func (tg *TgBot) processSubscribeCommand(ctx context.Context, msg tgbotapi.Message) error {
	return tg.repo.AddUser(ctx, msg.From.ID)
}

func (tg *TgBot) processUnsubscribeCommand(ctx context.Context, msg tgbotapi.Message) error {
	return tg.repo.RemoveUser(ctx, msg.From.ID)
}