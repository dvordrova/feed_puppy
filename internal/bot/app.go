package bot

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"github.com/dvordrova/feed_puppy/internal/database"
	"github.com/dvordrova/feed_puppy/internal/layout"
	"github.com/dvordrova/feed_puppy/internal/utils"

	_ "modernc.org/sqlite"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type App struct {
	layout   *layout.Layout
	bot      *tele.Bot
	handler  *Handler
	database *sql.DB
	wg       sync.WaitGroup
}

func CreateApp(config string) *App {
	lt, err := layout.New(config)
	if err != nil {
		log.Fatal("can't create layout: ", err)
	}
	bot, err := tele.NewBot(lt.Settings())
	if err != nil {
		log.Fatal("can't create bot: ", err)
	}

	ctx := context.Background()
	db, err := sql.Open("sqlite", utils.GetEnv("SQLITE_PATH", "database.sqlite3"))
	if err != nil {
		log.Fatal("can't open database: ", err)
	}
	db.SetMaxOpenConns(1)
	queries := database.New(db)
	handler := &Handler{layout: lt, db: db, queries: queries}

	bot.DeleteCommands()
	bot.SetCommands(lt.CommandsLocale("en"))
	bot.SetCommands(lt.CommandsLocale("ru"), "ru")
	bot.Use(middleware.Logger())
	bot.Use(lt.Middleware("ru", func(u *tele.User) (lang string) {

		user, err := queries.GetUser(ctx, u.ID)
		if err == nil {
			lang = user.Language
		} else if err == sql.ErrNoRows {
			lang = utils.GetLanguage(u.LanguageCode)
			queries.NewUser(ctx, database.NewUserParams{
				TelegramID: u.ID,
				Language:   lang,
				Name:       utils.GetName(u.FirstName, u.LastName),
			})
		}
		return
	}))

	bot.Handle("/start", handler.OnStart)
	bot.Handle("/help", handler.OnHelp)
	bot.Handle("/settings", handler.OnSettings)
	bot.Handle("/new_dog", handler.OnNewDog)
	bot.Handle("/cur_dog", handler.OnChangeCurDog)
	bot.Handle("/share_dog_for_owner", handler.OnShareDogForOwner)
	bot.Handle("/share_dog_for_reader", handler.OnShareDogForReader)
	bot.Handle("/unshare_all_dogs", handler.OnUnshareAllDogs)
	bot.Handle("/unsubscribe", handler.OnUnsubscribe)
	bot.Handle("/action", handler.OnActionCommand)

	bot.Handle(lt.Callback("lang"), handler.OnLang)
	bot.Handle(lt.Callback("action"), handler.OnAction)
	bot.Handle(tele.OnText, handler.OnText)

	return &App{layout: lt, bot: bot, database: db, handler: handler}
}

func (app *App) RunBot() {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		log.Println("Started bot")
		app.bot.Start()
	}()
}

func (app *App) RunProcesses() {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		// app.bot.Start()
	}()
}

func (app *App) Finish() {
	app.database.Close()
}

func (app *App) Wait() {
	app.wg.Wait()
}
