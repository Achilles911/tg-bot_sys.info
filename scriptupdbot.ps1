docker stop telegram-bot
docker remove telegram-bot
docker build -t telegram-bot .
docker run -d --name telegram-bot telegram-bot
