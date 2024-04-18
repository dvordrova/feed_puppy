package bot

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/dvordrova/feed_puppy/internal/database"
	"github.com/dvordrova/feed_puppy/internal/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) OnText(c tele.Context) (err error) {
	var (
		ctx          = context.TODO()
		tx           *sql.Tx
		qtx          *database.Queries
		user         database.User
		dogId        int64
		dogBirthDate time.Time
	)

	if tx, err = h.db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()
	qtx = h.queries.WithTx(tx)
	if user, err = qtx.GetUser(ctx, c.Sender().ID); err != nil {
		return
	}

	if user.State == string(utils.UserStateNewDog) {
		log.Println("new dog", c.Data())
		mapMaleFemale := map[string]string{
			"f":       "Female",
			"fem":     "Female",
			"female":  "Female",
			"ж":       "Female",
			"жен":     "Female",
			"женский": "Female",
			"m":       "Male",
			"male":    "Male",
			"м":       "Male",
			"муж":     "Male",
			"мужской": "Male",
		}
		dogParams := strings.Split(c.Message().Text, ",")
		if len(dogParams) != 4 {
			return c.Send(h.layout.Text(c, "register_new_dog"), tele.Silent)
		}
		caser := cases.Title(language.Und)
		dogNameStr := caser.String(strings.TrimSpace(dogParams[0]))
		dogBirthDateStr := strings.TrimSpace(dogParams[1])
		dogBreedStr := caser.String(strings.TrimSpace(dogParams[2]))
		dogSexStr := strings.ToLower(strings.TrimSpace(dogParams[3]))

		dogSex, ok := mapMaleFemale[dogSexStr]
		if !ok {
			return c.Send(h.layout.Text(c, "register_new_dog"), tele.Silent)
		}

		layout := "02.01.2006" // Day.Month.Year layout
		dogBirthDate, err = time.Parse(layout, dogBirthDateStr)
		if err != nil {
			return c.Send(h.layout.Text(c, "register_new_dog"), tele.Silent)
		}
		dogId, err = qtx.NewDog(ctx, database.NewDogParams{
			Name:      dogNameStr,
			BirthDate: dogBirthDateStr,
			Breed:     dogBreedStr,
			Sex:       dogSex,
		})
		qtx.NewSubscription(ctx, database.NewSubscriptionParams{
			UserID: user.ID,
			DogID:  dogId,
			Type:   "owner",
		})
		if err != nil {
			return err
		}
		err = qtx.SetUserCurDog(ctx, database.SetUserCurDogParams{
			ID:         user.ID,
			CurrentDog: sql.NullInt64{Int64: dogId, Valid: true},
		})
		if err != nil {
			return err
		}
		err = qtx.SetUserState(ctx, database.SetUserStateParams{
			ID:    user.ID,
			State: string(utils.UserStateDogSelected),
		})
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return
		}

		daysFromBirth := int64(time.Since(dogBirthDate).Hours() / 24)
		err = c.Send(
			h.layout.Text(c, "msg_new_dog_created", struct {
				DogName          string
				DogDaysFromBirth int64
			}{
				DogName:          dogNameStr,
				DogDaysFromBirth: daysFromBirth,
			}),
			tele.Silent,
		)
		if err != nil {
			return
		}
		return c.Send(
			h.layout.Text(c, "msg_choose_action"),
			h.layout.Markup(c, "action"),
			tele.Silent,
		)
	}
	log.Println(user.State)
	if user.State == string(utils.UserStateActionYum) ||
		user.State == string(utils.UserStateActionFeed) ||
		user.State == string(utils.UserStateActionWeigh) {
		addInfo := c.Text()
		qtx.NewAction(ctx, database.NewActionParams{
			UserID:         user.ID,
			DogID:          user.CurrentDog.Int64,
			ActionID:       int64(utils.ActionsMap[user.State]),
			AdditionalInfo: sql.NullString{String: addInfo, Valid: true},
			Timestamp:      time.Now().Unix(),
		})
		dog, err := qtx.GetDog(ctx, user.CurrentDog.Int64)
		if err != nil {
			return err
		}

		err = qtx.SetUserState(ctx, database.SetUserStateParams{
			ID:    user.ID,
			State: string(utils.UserStateDogSelected),
		})
		if err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
		utils.NotifyAllDogSubscribers(
			c.Bot(),
			ctx,
			h.queries,
			dog.ID,
			h.layout.Text(c, "msg_action_tmpl_add_info", struct {
				DateTime string
				DogName  string
				UserName string
				AddInfo  string
				Action   string
			}{
				DateTime: time.Now().Format("02.01.2006 15:04"),
				DogName:  dog.Name,
				UserName: c.Sender().FirstName,
				AddInfo:  addInfo,
				Action:   user.State,
			}),
		)
		return c.Send(
			h.layout.Text(c, "msg_choose_action"),
			h.layout.Markup(c, "action"),
			tele.Silent,
		)
	}

	if err = tx.Commit(); err != nil {
		return
	}
	// defer c.Bot().Delete(c.Message())
	return c.Send("idk how to respond")
}
