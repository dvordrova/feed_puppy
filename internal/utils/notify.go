package utils

import (
	"context"
	"log"

	"github.com/dvordrova/feed_puppy/internal/database"
	tele "gopkg.in/telebot.v3"
)

func NotifyAllDogSubscribers(bot *tele.Bot, ctx context.Context, q *database.Queries, dogId int64, notification string) {
	subscribers, err := q.SelectDogsSubscribers(ctx, dogId)
	if err != nil {
		return
	}

	for _, subscriber := range subscribers {
		if _, err := bot.Send(&tele.User{ID: subscriber.TelegramID}, notification); err != nil {
			log.Println("can't send notification to user: ", err)
		}
	}
}
