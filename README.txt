Flatbot is a Telegram notification tool for rightmove.co.uk URLs.


RUN

$ export FLATBOT_TELEGRAM_BOT_API_TOKEN=<...>
$ export FLATBOT_TELEGRAM_CHAT_ID=<...>
$ go build
$ ./flatbot -dry-run -once <URL>
$ ./flatbot <URL>


TEST

$ go test
$ python3 -m http.server testdata
$ ./flatbot -dry-run -once http://localhost:8000/2025-02-19-basic.html
