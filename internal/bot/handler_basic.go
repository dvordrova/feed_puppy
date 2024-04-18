package bot

import (
	"context"
	"database/sql"
	"strings"

	"github.com/dvordrova/feed_puppy/internal/database"
	"github.com/dvordrova/feed_puppy/internal/utils"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) OnSettings(c tele.Context) error {
	defer c.Bot().Delete(c.Message())
	return c.Send(
		h.layout.Text(c, "settings_msg"),
		h.layout.Markup(c, "settings"),
		tele.Silent,
	)
}

func (h *Handler) OnHelp(c tele.Context) error {
	defer c.Bot().Delete(c.Message())
	return c.Send(h.layout.Text(c, "msg_help"), tele.Silent)
}

func (h *Handler) OnStart(c tele.Context) (err error) {
	var (
		ctx  = context.TODO()
		cmd  = strings.Split(c.Message().Text, " ")
		tx   *sql.Tx
		qtx  *database.Queries
		user database.User
		dog  database.Dog
	)

	if tx, err = h.db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()
	qtx = h.queries.WithTx(tx)
	if user, err = qtx.GetUser(ctx, c.Sender().ID); err != nil {
		return
	}

	// defer c.Bot().Delete(c.Message())
	if len(cmd) == 1 {
		return c.Send(h.layout.Text(c, "start_without_dog"), tele.Silent)
	}

	invite, err := qtx.GetInvite(ctx, cmd[1])
	if err == sql.ErrNoRows {
		return c.Send(h.layout.Text(c, "start_without_dog"), tele.Silent)
	}

	qtx.MarkInviteUsed(ctx, invite.ID)
	qtx.ClearSubscriptions(ctx, database.ClearSubscriptionsParams{
		UserID: user.ID,
		DogID:  dog.ID,
		Type:   "reader",
	})

	if invite.Type == "owner" {
		qtx.ClearSubscriptions(ctx, database.ClearSubscriptionsParams{
			UserID: user.ID,
			DogID:  dog.ID,
			Type:   "owner",
		})
		qtx.NewSubscription(ctx, database.NewSubscriptionParams{
			DogID:  invite.DogID,
			UserID: user.ID,
			Type:   invite.Type,
		})
		qtx.SetUserCurDog(ctx, database.SetUserCurDogParams{
			ID:         user.ID,
			CurrentDog: sql.NullInt64{Int64: invite.DogID, Valid: true},
		})
		qtx.SetUserState(ctx, database.SetUserStateParams{
			ID:    user.ID,
			State: string(utils.UserStateDogSelected),
		})
		dog, err = qtx.GetDog(ctx, invite.DogID)
		if err != nil {
			return
		}
		if err = tx.Commit(); err != nil {
			return
		}
		c.Send(
			h.layout.Text(c, "start_with_dog_owner", struct {
				YourDog string
			}{
				YourDog: dog.Name,
			}),
			tele.Silent,
		)
		return c.Send(
			h.layout.Text(c, "msg_choose_action"),
			h.layout.Markup(c, "action"),
			tele.Silent,
		)
	}
	qtx.NewSubscription(ctx, database.NewSubscriptionParams{
		DogID:  invite.DogID,
		UserID: user.ID,
		Type:   invite.Type,
	})
	if err = tx.Commit(); err != nil {
		return
	}

	return c.Send(h.layout.Text(c, "start_with_dog_reader"))
}
