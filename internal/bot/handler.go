package bot

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/dvordrova/feed_puppy/internal/database"
	"github.com/dvordrova/feed_puppy/internal/layout"
	"github.com/dvordrova/feed_puppy/internal/utils"

	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	layout  *layout.Layout
	db      *sql.DB
	queries *database.Queries
}

func (h *Handler) OnLang(c tele.Context) error {
	log.Println("/callback[lang]")
	lang := c.Data()
	ctx := context.TODO()
	h.queries.SetUserLanguage(ctx, database.SetUserLanguageParams{
		TelegramID: c.Sender().ID,
		Language:   lang,
	})

	h.layout.SetLocale(c, lang)
	defer c.Bot().Delete(c.Message())
	return c.Bot().Respond(c.Callback(), &tele.CallbackResponse{Text: h.layout.Text(c, "lang_setted")})
}

func (h *Handler) OnAction(c tele.Context) (err error) {
	var (
		ctx    = context.TODO()
		action = c.Data()
		tx     *sql.Tx
		qtx    *database.Queries
		user   database.User
		dog    database.Dog
	)

	if tx, err = h.db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()
	qtx = h.queries.WithTx(tx)
	if user, err = qtx.GetUser(ctx, c.Sender().ID); err != nil {
		return
	}
	if user.State != string(utils.UserStateDogSelected) || !user.CurrentDog.Valid {
		c.Bot().Respond(c.Callback(), &tele.CallbackResponse{Text: h.layout.Text(c, "select_dog_first")})
		err = errors.New("no selected dog")
		return
	}

	if action == "action_feed" || action == "action_yum" || action == "action_weigh" {
		qtx.SetUserState(ctx, database.SetUserStateParams{
			ID:    user.ID,
			State: action,
		})
		if err = tx.Commit(); err != nil {
			return
		}
		err = c.Edit(h.layout.Text(c, action+"_text_needed"))
	} else {
		dog, err = qtx.GetDog(ctx, user.CurrentDog.Int64)
		qtx.NewAction(ctx, database.NewActionParams{
			UserID:    user.ID,
			DogID:     dog.ID,
			ActionID:  int64(utils.ActionsMap[action]),
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			return
		}
		utils.NotifyAllDogSubscribers(c.Bot(), ctx, h.queries, dog.ID,
			h.layout.Text(c, "msg_action_tmpl", struct {
				DateTime string
				DogName  string
				UserName string
				Action   string
			}{
				DateTime: time.Now().Format("02.01.2006 15:04"),
				DogName:  dog.Name,
				UserName: c.Sender().FirstName,
				Action:   action,
			}))
		defer c.Bot().Delete(c.Message())
		err = c.Send(
			h.layout.Text(c, "msg_choose_action"),
			h.layout.Markup(c, "action"),
			tele.Silent,
		)
	}
	return
}
