Flatbot is a Telegram notification tool for rightmove.co.uk URLs.


RUN

$ export FLATBOT_TELEGRAM_BOT_API_TOKEN=<...>
$ export FLATBOT_TELEGRAM_CHAT_ID=<...>
$ go build && ./flatbot <URL>


DECRYPT SECRETS

$ . <(gpg --decrypt --quiet misc/telegram-secrets.sh.asc)
