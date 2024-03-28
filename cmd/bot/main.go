package main

import "github.com/dvordrova/feed_puppy/internal/bot"

func main() {
	app := bot.CreateApp("bot.yml")
	defer app.Finish()

	app.RunBot()
	app.Wait()
}
