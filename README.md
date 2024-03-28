Telegram bot for tracking your dog activity in a very simple way

Storage - sqlite

# Fast start

1. Create a new bot with BotFather
2. Copy the token
3. Save token as a temporary env for start

```bash
read -s TELEGRAM_TOKEN
```

4. install air for fast reload of your app

```bash
go install github.com/cosmtrek/air@latest
```

5. [https://docs.sqlc.dev/en/latest/overview/install.html](install) sqlc for generating code from sql

6. Run migrations and start the bot with air help

```bash
TELEGRAM_TOKEN=$TELEGRAM_TOKEN air
```

6. Go to your bot and type /start

7. check your base with sqlite3

```bash
sqlite3 -readonly database.sqlite3
```

8. deploy it whenever you want
