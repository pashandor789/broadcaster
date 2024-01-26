package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	subscribeCommand   = "subscribe"
	unsubscribeCommand = "unsubscribe"
	startCommand       = "start"
)

type repository interface {
	AddUser(ctx context.Context, ID int64) error
	RemoveUser(ctx context.Context, ID int64) error
	GetUsersID(ctx context.Context) ([]int64, error)
}
type TgBot struct {
	repo repository

	*tgbotapi.BotAPI
}

func NewTgBot(cfg BotConfig, repo repository) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		return nil, err
	}

	return &TgBot{
		BotAPI: bot,
		repo:   repo,
	}, nil
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
	case startCommand:
		err = tg.processStartCommand(ctx, msg)
	case subscribeCommand:
		err = tg.processSubscribeCommand(ctx, msg)
	case unsubscribeCommand:
		err = tg.processUnsubscribeCommand(ctx, msg)
	default:
		m := tgbotapi.NewMessage(msg.Chat.ID, "No such command.")
		_, _ = tg.BotAPI.Send(m)
	}

	if err != nil {
		log.Printf("error occured: %v", err)

		errMsg := tgbotapi.NewMessage(msg.Chat.ID, "Something went wrong. ٩(͡๏̯͡๏)۶")
		_, _ = tg.BotAPI.Send(errMsg)
	}
}

func (tg *TgBot) processStartCommand(ctx context.Context, msg tgbotapi.Message) error {
	m := tgbotapi.NewMessage(msg.From.ID, "/subscribe, /unsubscribe")
	_, err := tg.BotAPI.Send(m)

	return err
}

func (tg *TgBot) processSubscribeCommand(ctx context.Context, msg tgbotapi.Message) error {
	err := tg.repo.AddUser(ctx, msg.From.ID)

	if err != nil {
		return err
	}

	m := tgbotapi.NewMessage(msg.From.ID, "You have successfully subscribed to junk!")
	_, err = tg.BotAPI.Send(m)

	if err != nil {
		return err
	}

	return nil
}

func (tg *TgBot) processUnsubscribeCommand(ctx context.Context, msg tgbotapi.Message) error {
	err := tg.repo.RemoveUser(ctx, msg.From.ID)

	if err != nil {
		return err
	}

	m := tgbotapi.NewMessage(msg.From.ID, "You have successfully unsubscribed from junk!")
	_, err = tg.BotAPI.Send(m)

	if err != nil {
		return err
	}

	return nil
}

func (tg *TgBot) BroadcastSubscribers(ctx context.Context, content string) error {
	ids, err := tg.repo.GetUsersID(ctx)

	if err != nil {
		return err
	}

	m := tgbotapi.NewMessage(-1, content)

	for _, id := range ids {
		m.ChatID = id
		_, _ = tg.Send(m)
	}

	return nil
}
