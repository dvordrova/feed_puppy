Telegram bot for tracking your dog activity in a very simple way

Storage - sqlite

# Fast start

1. Use golang 1.22
2. Create a new bot with [BotFather](https://t.me/BotFather)
3. Save token as a temporary env for start. Copy paste next command:

```bash
read -s TELEGRAM_TOKEN
```

4. Copy the bot's token into the input field of the previous command once and press enter (pasted token will be invisible)

5. [install](https://docs.sqlc.dev/en/latest/overview/install.html) sqlc for generating code from sql. For macos brew users just copy paste next command:

```bash
brew install sqlc
```

6. Run migrations and start the bot with air help. Just copy paste next command:

```bash
TELEGRAM_TOKEN=$TELEGRAM_TOKEN go run github.com/cosmtrek/air@latest
```

7. Go to your bot and type /start

8. check your base with sqlite3

```bash
sqlite3 -readonly database.sqlite3
```

9. deploy it whereever you want
